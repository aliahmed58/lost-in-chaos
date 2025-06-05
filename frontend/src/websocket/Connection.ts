// websocket connection

import { GameEngine } from "../engine/engine"

export const connect = (engine: GameEngine) => {
    let socket = new WebSocket(`ws://192.168.18.196/websocket?uuid=${engine.mainPlayer.uuid}`)

    socket.onopen = e => {
        console.log('connected to server successfully')
    }

    socket.onclose = e => {
        engine.players.clear()
    }

    socket.onmessage = e => {
        console.log(e.data)
    }

    socket.onerror = e => {
        console.log(e)
    }
}