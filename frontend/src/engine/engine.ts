import { Application, ApplicationOptions, Assets, DestroyOptions, RendererDestroyOptions, Texture } from "pixi.js";
import { getResolution } from "./utils/getResoultion";
import { Player } from "../entities/Player";
import { connect } from "../websocket/Connection";

export class GameEngine extends Application {

    public players: Map<string, Player> = new Map<string, Player>();
    public mainPlayer!: Player;
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
        connect(this)
    }

    public override destroy(rendererDestroyOptions?: RendererDestroyOptions, options?: DestroyOptions): void {
        super.destroy(rendererDestroyOptions, options);
    }
}