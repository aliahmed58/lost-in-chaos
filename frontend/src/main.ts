import { Application, Assets, Sprite } from "pixi.js";
import { GameEngine } from "./engine/engine";
import { Player } from "./entities/Player";

const engine = new GameEngine();
(async () => {

	// Initialize the application
	await engine.init({ background: "#1099bb", resizeTo: window });

	engine.connectToServer()

	// Listen for animate update
	engine.ticker.add((time) => {
		engine.update(time.deltaTime)
	});
})();
