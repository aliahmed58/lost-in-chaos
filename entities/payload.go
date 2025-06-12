package entities

type Payload struct {
	Type string `json:"type"`
}

type JoinNotif struct {
	Type string `json:"type"`
	Uuid string `json:"uuid"`
}
