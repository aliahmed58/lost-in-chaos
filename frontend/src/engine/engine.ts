import { Application, ApplicationOptions, DestroyOptions, RendererDestroyOptions } from "pixi.js";
import { getResolution } from "./utils/getResoultion";

export class GameEngine extends Application {
    public override async init(options: Partial<ApplicationOptions>): Promise<void> {
        options.resizeTo ??= window;
        options.resolution ??= getResolution();

        await super.init(options);
        
        // Append the application canvas to the document body
        document.getElementById("pixi-container")!.appendChild(this.canvas);
    }

    public override destroy(rendererDestroyOptions?: RendererDestroyOptions, options?: DestroyOptions): void {
        super.destroy(rendererDestroyOptions, options);
    }
}