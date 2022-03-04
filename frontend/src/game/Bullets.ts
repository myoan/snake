import 'phaser';
import TextureKey from '../enum/TextureKey';

export class Bullet extends Phaser.Physics.Arcade.Sprite {
  scene: Phaser.Scene;
  ttl: Phaser.Time.TimerEvent;

  constructor(scene: Phaser.Scene, x: number, y: number) {
    super(scene, x, y, TextureKey.Bullet);
    this.scene = scene;
  }

  fire(x: number, y: number, vx: number, vy: number) {
    this.body.reset(x, y);

    this.scene.physics.add.existing(this);

    this.setActive(true);
    this.setVisible(true);

    this.setVelocity(vx, vy);
    this.ttl = this.scene.time.addEvent({delay: 500, callback: this.dead, args: [this]});
  }

  dead(args) {
    const bullet: Bullet = args;
    bullet.destroy();
  }

  preUpdate(time, delta) {
    super.preUpdate(time, delta);

    if (this.y < 0 || this.x < 0 || this.x > 2160, this.y > 2160) {
      this.setActive(false);
      this.setVisible(false);
    }
  }
}

export class Bullets extends Phaser.Physics.Arcade.Group {
  scene: Phaser.Scene;

  constructor(scene: Phaser.Scene) {
    super(scene.physics.world, scene);
    this.scene = scene;

    this.createMultiple({
      frameQuantity: 5,
      key: TextureKey.Bullet,
      active: false,
      visible: false,
      classType: Bullet,
    });
  }

  fireBullet(x: number, y: number, vx: number, vy: number) {
    const bullet = this.getFirstDead(true);
    this.scene.physics.add.existing(bullet);
    bullet.fire(x, y, vx, vy);
  }
}

export default Bullets;
