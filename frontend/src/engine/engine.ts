import { Application, ApplicationOptions, DestroyOptions, RendererDestroyOptions } from "pixi.js";
import { getResolution } from "./utils/getResoultion";
import { Player } from "../entities/Player";

export class GameEngine extends Application {

    public players: Map<string, Player> = new Map<string, Player>();

    public override async init(options: Partial<ApplicationOptions>): Promise<void> {
        options.resizeTo ??= window;
        options.resolution ??= getResolution();

        await super.init(options);
        
        // Append the application canvas to the document body
        document.getElementById("pixi-container")!.appendChild(this.canvas);
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

    public override destroy(rendererDestroyOptions?: RendererDestroyOptions, options?: DestroyOptions): void {
        super.destroy(rendererDestroyOptions, options);
    }
}