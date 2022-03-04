import 'phaser';
import Ship from './Ship';

export default class Player extends Ship {
  scene: Phaser.Scene;

  constructor(scene: Phaser.Scene, id: string, x: number, y: number) {
    super(scene, id, x, y);

    this.scene = scene;
    scene.add.existing(this);
    scene.physics.world.enable(this);
  }

  setInputSources() {
    this.scene.input.keyboard.on('keydown-W', () => this.moveUpper(), this.scene);
    this.scene.input.keyboard.on('keydown-S', () => this.moveDowner(), this.scene);
    this.scene.input.keyboard.on('keydown-A', () => this.moveLeft(), this.scene);
    this.scene.input.keyboard.on('keydown-D', () => this.moveRight(), this.scene);
    this.scene.input.keyboard.on('keyup', () => this.stop(), this.scene);
    this.scene.input.on('pointerdown', () => this.fire() );
    this.scene.input.on('pointermove', (ptr: Phaser.Input.Pointer) => {
      this.setDirection(ptr.worldX, ptr.worldY);
    });
  }

  setOverlap(enemies: Ship[]) {
    this.scene.physics.add.overlap(enemies, this.bullets, this.hit);
  }
}
