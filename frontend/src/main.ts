import { Application, Assets, Sprite } from "pixi.js";
import { GameEngine } from "./engine/engine";
import { Player } from "./entities/Player";

const engine = new GameEngine();
(async () => {

	// Initialize the application
	await engine.init({ background: "#1099bb", resizeTo: window });

	// Load the bunny texture
	const texture = await Assets.load("/assets/bunny.png");

	const player = new Player(texture);
	player.anchor.set(0.5);
	player.position.set(engine.screen.width / 2, engine.screen.height / 2);
	player.addListeners()
	engine.addPlayer(player)

	// Listen for animate update
	engine.ticker.add((time) => {
		player.move(time.deltaTime)
	});
})();
