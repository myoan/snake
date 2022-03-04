import 'phaser';
import Coordinate from '../lib/Coordinate';
import Vector from '../lib/Vector';
import TextureKey from '../enum/TextureKey';
import {Bullet, Bullets} from './Bullets';

const SHIP_SPEED = 200;
const SHIP_HP = 3;
const BULLET_SPEED = 500;

export default class Ship extends Phaser.GameObjects.Container {
  scene: Phaser.Scene;
  cood: Coordinate;
  v: number;
  bullets: Bullets;
  id: string;
  hp: number;

  constructor(scene: Phaser.Scene, id: string, x: number, y: number) {
    super(scene, x, y);
    this.scene = scene;
    this.id = id;
    this.hp = SHIP_HP;
    this.cood = new Coordinate(new Vector(x, y), 0);

    const obj = this.scene.physics.add.image(0, 0, TextureKey.Ship).setOrigin(0.5, 0.5);
    this.add(obj);

    // physicsにaddしないとsetVelocityなどが利用できない
    this.scene.physics.add.existing(this);

    const body = this.body as Phaser.Physics.Arcade.Body;
    this.setSize(40, 40);
    body.setCircle(20);
    body.setVelocity(0, 0);
    body.setCollideWorldBounds(true);
    body.setBounce(1, 1);
    this.rotation = this.radian(90);

    this.bullets = new Bullets(scene);
  }

  pos() {
    return this.cood.convertToWorld(new Vector(0, 0));
  }

  rotate(x: number) {
    this.rotation += this.radian(x);
    this.cood.rotate(x);
  }

  setDirection(x: number, y: number) {
    const local = this.cood.convertToLocal(new Vector(x, y));

    if (local.y > 0) {
      if (local.y < 5) {
        this.rotate(1);
      } else {
        this.rotate(5);
      }
    } else {
      if (local.y > -5) {
        this.rotate(-1);
      } else {
        this.rotate(-5);
      }
    }
  }

  setDirectionPhi(phi: number) {
    const diff = phi - this.cood.theta;
    this.rotate(diff);
  }

  moveUpper() {
    this.move(SHIP_SPEED, 0);
  }

  moveDowner() {
    this.move(SHIP_SPEED, -180);
  }

  moveLeft() {
    this.move(SHIP_SPEED, -90);
  }

  moveRight() {
    this.move(SHIP_SPEED, 90);
  }

  stop() {
    this.move(0, 0);
  }

  fire() {
    if (!this.active) return;

    const pos = this.pos();
    const v = this.cood.directionToWorld(0, BULLET_SPEED);
    this.bullets.fireBullet(pos.x, pos.y, v.x, v.y);
  }

  hit(obj1: Phaser.GameObjects.GameObject, obj2: Phaser.GameObjects.GameObject) {
    const ship = obj1 as Ship;
    ship.hp -= 1;
    if (ship.hp < 0) {
      ship.destroy();
    }
    const bullet = obj2 as Bullet;
    bullet.destroy();
  }

  private move(r: number, phai: number) {
    const v = this.cood.directionToWorld(phai, r);
    const body = this.body as Phaser.Physics.Arcade.Body;
    if (body == undefined) return;
    body.setVelocity(v.x, v.y);
    this.updatePos();
  }

  private updatePos() {
    this.cood.pos.x = this.x;
    this.cood.pos.y = this.y;
  }

  private radian(deg: number): number {
    return deg * Math.PI / 180;
  }
}
