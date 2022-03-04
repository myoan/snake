import 'phaser';
import TextureKey from '../enum/TextureKey';

export default class Preloader extends Phaser.Scene {
  graphics: Phaser.GameObjects.Graphics;
  conn: WebSocket;

  constructor() {
    super('preloader');
  }

  preload() {
    this.graphics = this.add.graphics();
    this.definePlayerTexture();
    this.defineBullet();
    this.conn = new WebSocket('ws://' + document.location.hostname + ':8008/ws');
  }

  private definePlayerTexture() {
    this.graphics.fillStyle(0x00fd00, 1.0);
    this.graphics.fillCircle(30, 30, 20);

    this.graphics.lineStyle(1, 0x00fd00, 1);
    this.graphics.beginPath();
    this.graphics.moveTo(30, 0);
    this.graphics.lineTo(24, 6);
    this.graphics.lineTo(36, 6);
    this.graphics.closePath();
    this.graphics.fillPath();
    this.graphics.generateTexture(TextureKey.Ship, 60, 60);
    this.graphics.clear();
  }

  defineBullet() {
    this.graphics.fillStyle(0x00fd00, 1.0);
    this.graphics.fillCircle(2, 2, 2);
    this.graphics.generateTexture(TextureKey.Bullet, 4, 4);
    this.graphics.clear();
  }

  drawMap(width: number, height: number) {
    // draw Outline
    this.graphics.lineStyle(5, 0x00fd00, 1);
    this.graphics.beginPath();
    this.graphics.moveTo(0, 0);
    this.graphics.lineTo(0, height);
    this.graphics.lineTo(width, height);
    this.graphics.lineTo(width, 0);
    this.graphics.closePath();
    this.graphics.strokePath();

    // draw Mesh
    this.graphics.lineStyle(1, 0x00fd00, 1);
    this.graphics.beginPath();

    for (let i = 1; i <= 10; i++) {
      this.graphics.moveTo(width * i / 10, 0);
      this.graphics.lineTo(width * i / 10, height);
    }
    for (let i = 1; i <= 10; i++) {
      this.graphics.moveTo(0, height * i / 10);
      this.graphics.lineTo(width, height * i / 10);
    }
    this.graphics.strokePath();
  }

  create() {
    this.scene.start('game', this.conn);
  }
}
