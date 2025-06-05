import { Sprite, SpriteOptions, Texture } from "pixi.js";

// movement consts
const
    ACCELERATION = 1,
    FRICTION = 0.8;

export class Player extends Sprite {

    private _vx: number = 0;
    private _vy: number = 0;
    private _uuid: string = crypto.randomUUID();

    constructor(options?: SpriteOptions | Texture) {
        super(options)
    }

    public handleInput(event: KeyboardEvent): void {
        switch (event.key) {
            case "a":
				this._vx -= ACCELERATION;
				break;
			case "w":
				this._vy -= ACCELERATION;
				break;
			case "s":
				this._vy += ACCELERATION;
				break;
			case "d":
				this._vx += ACCELERATION;
				break;
        }
    }

    public addListeners(): void {
        window.addEventListener('keydown', e => {
            this.handleInput(e)
        })
    }

    public move(deltaTime: number): void {
        super.x += this._vx * deltaTime
        super.y += this._vy * deltaTime 

        this._vx *= FRICTION
        this._vy *= FRICTION

		if (Math.abs(this._vx) < 0.01) this._vx = 0;
		if (Math.abs(this._vy) < 0.01) this._vy = 0;
    }

    public get vy(): number {
        return this._vy
    }
    
    public get vx(): number {
        return this._vx
    }

    public get uuid(): string {
        return this._uuid;
    }
}