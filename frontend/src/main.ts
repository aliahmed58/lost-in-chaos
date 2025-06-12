import { GameEngine } from "./engine/engine";

const engine = new GameEngine();
(async () => {

	// Initialize the application
	await engine.init({ background: "#1099bb", resizeTo: window });

	console.log(engine.players.values())
	// Listen for animate update
	if (engine.socket.readyState === engine.socket.OPEN) {
		engine.socket.send(JSON.stringify({type: "player_join"}))
	}
	engine.ticker.add((time) => {
		engine.update(time.deltaTime)
	});
})();
