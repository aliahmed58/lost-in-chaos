import { Application, Assets, Sprite } from "pixi.js";
import { GameEngine } from "./engine/engine";

const engine = new GameEngine();
(async () => {

	interface Player {
		uuid: string;
		x: number;
		y: number;
	}
	const players: Record<string, Sprite> = {};
	const curr = {} as Player
	curr.uuid = crypto.randomUUID();
	const socket = new WebSocket(`ws://192.168.18.196/websocket?uuid=${curr["uuid"]}`);

	socket.onopen = function () {
		console.log(socket)
		console.log("Connected to WebSocket server.");
	};

	// Initialize the application
	await engine.init({ background: "#1099bb", resizeTo: window });

	// Load the bunny texture
	const texture = await Assets.load("/assets/bunny.png");

	// Create a bunny Sprite
	const bunny = new Sprite(texture);

	// Center the sprite's anchor point
	bunny.anchor.set(0.5);

	// Move the sprite to the center of the screen
	bunny.position.set(engine.screen.width / 2, engine.screen.height / 2);


	socket.onmessage = event => {
		console.log(event.data)
		let json = JSON.parse(event.data)
		if (!players[json.uuid]) {
			// add new bunny lmao
			const b = new Sprite(texture)
			b.position.set(json.x, json.y)
			engine.stage.addChild(b)
			players[json.uuid] = b
		}
		players[json.uuid].x = json.x
		players[json.uuid].y = json.y
	}

	// Add the bunny to the stage
	engine.stage.addChild(bunny);
	const friction = 0.7;
	const acc = 1;
	let vx = 0;
	let vy = 0;


	window.addEventListener('keydown', e => {
		switch (e.key) {
			case "a":
				vx -= acc;
				break;
			case "w":
				vy -= acc;
				break;
			case "s":
				vy += acc;
				break;
			case "d":
				vx += acc;
				break;
		}
		if (socket.readyState == socket.OPEN) {
			socket.send(JSON.stringify({uuid: curr.uuid, x: bunny.x, y: bunny.y}))
		}
	})

	// Listen for animate update
	engine.ticker.add((time) => {
		bunny.x += vx * time.deltaTime;
		bunny.y += vy * time.deltaTime;

		vx *= friction;
		vy *= friction;

		if (Math.abs(vx) < 0.01) vx = 0;
		if (Math.abs(vy) < 0.01) vy = 0;
	});
})();
