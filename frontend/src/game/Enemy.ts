import Coordinate from '../lib/Coordinate';
import 'phaser';
import Ship from './Ship';
// import Vector from 'lib/Vector';

/*
class Action {
  t: number;
  position: Vector;
  direction: number;
  move: number;
  shot: Boolean;

  constructor(t: number, direction: number, x: number, y: number) {
    this.t = t;
    this.position = new Vector(x, y);
    this.direction = direction;
  }
}

class ActionList {
  list: Action[];

  fetch(t: Number): Action {
    return this.list[0];
  }

  registerAction(t: number, direction: number, x: number, y: number) {
    const action = new Action(t, direction, x, y);
  }
}
*/

export default class Enemy extends Ship {
  scene: Phaser.Scene;
  i: number;
  // actionList: ActionList

  constructor(scene: Phaser.Scene, id: string, x: number, y: number) {
    super(scene, id, x, y);

    this.scene = scene;
    scene.add.existing(this);
    scene.physics.world.enable(this);
    this.i = -1;
  }

  randomWalk() {
    if (!this.active) return;

    const x = Phaser.Math.Between(0, 1000);
    const y = Phaser.Math.Between(0, 1000);
    this.setDirection(x, y);
    this.moveUpper();
  }

  randomShoot() {
    if (this.i == -1) {
      this.i = Phaser.Math.Between(0, 100);
    }
    this.i--;
    if (this.i < 0) {
      this.fire();
    }
  }

  setOverlap(enemies: Ship[]) {
    this.scene.physics.add.overlap(enemies, this.bullets, this.hit);
  }

  setPos(x: number, y: number) {
    this.cood.pos.x = x;
    this.cood.pos.y = y;
    this.setPosition(x, y);
  }

  interpolateLag() {
    this.moveUpper();
  }

  interpolateLag2(t: number, direction: number, x: number, y: number) {
    const frame = this.frame(t);
    console.log(`lag: ${t}, frame: ${frame}`);
    const interpolate = new Coordinate(this.cood.pos, this.cood.theta);
    for (let i = 0; i < frame; i++) {
      interpolate.move(200, 0);
    }
  }

  frame(t: number): number {
    const ft = 1000/60;
    const ret = Math.round(t / ft);
    if (ret == 0) return 1;
    return ret;
  }

  registerAction(t: Number, direction: Number, x: Number, y: Number) {
    // this.actionList.registerAction(t, direction, x, y);
  }
}
