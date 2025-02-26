package server

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

const (
	MagicVal            = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	SecWebSocketKey     = "Sec-WebSocket-Key"
	SecWebSocketAccept  = "Sec-WebSocket-Accept"
	SecWebSocketVersion = "Sec-WebSocket-Version"
	Upgrade             = "Upgrade"
	Connection          = "Connection"
)

var clientHeaders = []string{
	Upgrade,
	Connection,
	SecWebSocketKey,
	SecWebSocketVersion,
}

type UpgradeError struct {
	Err        error
	StatusCode int
}

func (uErr *UpgradeError) Error() string {
	return fmt.Sprintf("status: %d - err %v", uErr.StatusCode, uErr.Err)
}

func getClientHeaders(r *http.Request) (map[string]string, error) {
	c := make(map[string]string)
	for _, header := range clientHeaders {
		c[header] = r.Header.Get(header)
		if c[header] == "" {
			return nil, errors.New(header + " found empty in request. Cannot upgrade")
		}
		// TODO: if any header is wrong or misunderstood
	}
	return c, nil
}

type DataFrame struct {
	Fin           byte
	Opcode        byte
	Mask          byte
	MaskingKey    [4]byte
	PayloadLength uint
	Payload       []byte
}

func (d DataFrame) String() string {
	return fmt.Sprintf("FIN: %d | Mask: %d | Opcode: %d | Length: %d | Masking Key Len: %d\nMasking Key: %b\n",
		d.Fin, d.Mask, d.Opcode, d.PayloadLength, len(d.MaskingKey), d.MaskingKey)
}

func NewDataFrame(reader *bufio.Reader) (DataFrame, error) {
	frame := DataFrame{}
	byte1, err := reader.ReadByte()
	if err != nil {
		return frame, err
	}
	frame.Fin = (byte1 >> 7) & 1
	// fmt.Printf("%08b ", byte1)
	// fmt.Printf("fin byte: %08b %b\n", byte1, frame.Fin)
	frame.Opcode = byte1 & 0b00001111

	byte2, err := reader.ReadByte()
	if err != nil {
		return frame, err
	}
	frame.Mask = (byte2 >> 7) & 1
	// fmt.Printf("%08b ", byte2)
	// fmt.Printf("mask byte: %08b %b\n", byte2, frame.Mask)
	length, err := decodePayloadLength(byte2, reader)
	if err != nil {
		return frame, err
	}
	frame.PayloadLength = length

	maskingKey := make([]byte, 4)
	_, err = io.ReadFull(reader, maskingKey)
	if err != nil {
		return frame, err
	}
	frame.MaskingKey = [4]byte(maskingKey)
	// fmt.Printf("%08b ", maskingKey)

	frame.Payload = make([]byte, frame.PayloadLength)
	_, err = io.ReadFull(reader, frame.Payload)
	if err != nil {
		return frame, err
	}
	if frame.PayloadLength == 19376 {
		// fmt.Println(frame.Payload)
	}
	// fmt.Printf("%d\n\n", frame.PayloadLength)

	return frame, nil
}

func (df DataFrame) UnmaskData() []byte {
	if df.Mask == 1 {
		decodedData := make([]byte, df.PayloadLength)
		for i, d := range df.Payload {
			decodedData[i] = d ^ df.MaskingKey[i%4]
		}

		return decodedData
	}
	return df.Payload
}

type WebsocketConn struct {
	tcpConn net.Conn
}

func UpgradeConn(w http.ResponseWriter, r *http.Request) (*WebsocketConn, error) {

	headers, err := getClientHeaders(r)
	if err != nil {
		return nil, &UpgradeError{StatusCode: http.StatusBadRequest, Err: err}
	}

	// hijack connection to get the underlying tcp connection
	conn, _, err := http.NewResponseController(w).Hijack()

	if err != nil {
		return nil, err
	}

	secWsAccept := generateAcceptKey(headers[SecWebSocketKey])

	responseHeaders := "HTTP/1.1 101 Switching Protocols\r\n" +
		Upgrade + ": websocket\r\n" +
		Connection + ": Upgrade\r\n" +
		SecWebSocketAccept + ": " + secWsAccept + "\r\n\r\n"

	_, err = conn.Write([]byte(responseHeaders))
	if err != nil {
		return nil, err
	}

	return &WebsocketConn{tcpConn: conn}, nil
}

func (wConn *WebsocketConn) ReadMsg() {
	defer wConn.tcpConn.Close()
	reader := bufio.NewReader(wConn.tcpConn)

	for {
		dataFrame, err := NewDataFrame(reader)
		if dataFrame.Opcode != 1 {
			continue
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		// fmt.Println(dataFrame.String())

		// _ = dataFrame.UnmaskData()
		// fmt.Println(dataFrame.String())
		decodedData := dataFrame.UnmaskData()
		fmt.Println(string(decodedData))
		// fmt.Println(len(string(decodedData)))
	}
}

// generate the value for `Sec-WebSocket-Accept` that'll be sent to the client
func generateAcceptKey(key string) string {
	key += MagicVal
	hasher := sha1.New()
	hasher.Write([]byte(key))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func decodePayloadLength(byte2 byte, reader *bufio.Reader) (uint, error) {
	// Read bits 9-15 (inclusive) and interpret that as an unsigned integer. If it's 125 or less, then that's the length; you're done.
	length := uint(byte2 & 0b01111111)

	if length <= 125 {
		return length, nil
	}

	// If it's 126, Read the next 16 bits and interpret those as an unsigned integer. You're done
	if length == 126 {
		var byte3and4 [2]byte
		_, err := reader.Read(byte3and4[:])
		if err != nil {
			return 0, err
		}
		// fmt.Printf("%08b ", byte3and4)
		length = uint(binary.BigEndian.Uint16(byte3and4[:]))
	} else if length == 127 {
		// Read the next 64 bits and interpret those as an unsigned integer. (The most significant bit must be 0.) You're done.
		bytes3to8 := make([]byte, 8)
		_, err := reader.Read(bytes3to8)
		if err != nil {
			return 0, err
		}
		// fmt.Printf("%08b ", bytes3to8)
		length = uint(binary.BigEndian.Uint64(bytes3to8))

		if length&(1<<63) != 0 {
			return 0, errors.New("invalid payload length: msb should be 0")
		}
	}

	return length, nil

}
