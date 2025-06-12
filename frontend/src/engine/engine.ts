import { Application, ApplicationOptions, Assets, DestroyOptions, RendererDestroyOptions, Texture } from "pixi.js";
import { getResolution } from "./utils/getResoultion";
import { Player } from "../entities/Player";
import { connect } from "../websocket/Connection";

export class GameEngine extends Application {

    public players: Map<string, Player> = new Map<string, Player>();
    public mainPlayer!: Player;
    public socket!: WebSocket;
    private texture!: Texture;

    public override async init(options: Partial<ApplicationOptions>): Promise<void> {
        options.resizeTo ??= window;
        options.resolution ??= getResolution();

        await super.init(options);

        // Append the application canvas to the document body
        document.getElementById("pixi-container")!.appendChild(this.canvas);
        this.texture = await Assets.load("/assets/bunny.png")
        this.mainPlayer = new Player(this.texture)
        this.mainPlayer.anchor.set(0.5)
        this.mainPlayer.position.set(this.screen.width / 2, this.screen.height / 2);
        this.mainPlayer.addListeners()
        this.addPlayer(this.mainPlayer)
        this.connectToServer()
    }

    // add player to stage and map
    public addPlayer(p: Player): void {
        if (this.players.has(p.uuid)) {
            return
        }
        this.players.set(p.uuid, p)
        this.stage.addChild(p)
    }

    public getPlayer(uuid: string): Player {
        let player = this.players.get(uuid)
        if (!player) {
            throw new Error(`Player with ${uuid} is not present.`)
        }
        return player
    }

    public removePlayer(uuid: string): void {
        let player = this.players.get(uuid)
        if (!player) {
            return
        }
        console.log('removing player')
        this.players.delete(uuid)
        this.stage.removeChild(player)
    }

    public getPlayers(): Map<string, Player> {
        return this.players
    }

    // main game loop updates
    public update(deltaTime: number): void {
        this.mainPlayer.move(deltaTime)
    }

    public connectToServer(): void {
        this.socket = new WebSocket(`ws://192.168.18.196/websocket?uuid=${this.mainPlayer.uuid}`)

        this.socket.onopen = e => {
            console.log('connected to server successfully')
        }

        this.socket.onclose = e => {
        }

        this.socket.onmessage = e => {
            console.log(e.data)
            let data = JSON.parse(e.data)
            if (data.type === "new_join") {
                let p : Player = new Player(this.texture)
                p.position.set(this.screen.width / 2, this.screen.height / 2)
                p.uuid = data.uuid
                this.addPlayer(p)
            }
            if (data.type === "left") {
                let uuid: string = data.uuid
                this.removePlayer(uuid)
            }
        }

        this.socket.onerror = e => {
            console.log(e)
        }
    }

    public override destroy(rendererDestroyOptions?: RendererDestroyOptions, options?: DestroyOptions): void {
        super.destroy(rendererDestroyOptions, options);
    }
}