import Ship from '../game/Ship';
import Player from '../game/Player';
import Enemy from '../game/Enemy';
import Bullets from '../game/Bullets';

const SCREEN_WIDTH = 1200;
const SCREEN_HEIGHT = 1000;
const PLAYER_NUM = 2;

let text: Phaser.GameObjects.Text;
const sTimes = Array<number>(PLAYER_NUM);
export default class Game extends Phaser.Scene {
  g: Phaser.GameObjects.Graphics;
  player: Player;
  enemies: Array<Enemy>;
  bullets: Bullets;
  conn: WebSocket;

  constructor(conn: WebSocket) {
    super('game');
  }

  create(conn: WebSocket) {
    this.cameras.main.setBounds(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT);
    this.physics.world.setBounds(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT);

    this.createMap(SCREEN_WIDTH, SCREEN_HEIGHT);

    this.player = new Player(this, 'player', 400, 300);
    this.player.setInputSources();

    text = this.add.text(10, 10, '', {font: '16px Courier', color: '#fdfdfd'}).setScrollFactor(0);
    this.createEnemies(PLAYER_NUM - 1);

    this.cameras.main.startFollow(this.player, true, 0.5, 0.5);
    this.conn = conn;

    // initialize sTimes
    for (let i = 0; i < PLAYER_NUM; i++) {
      sTimes[i] = -1;
    }

    conn.onmessage = (event) => {
      console.log(`response: ${event.data}`);
      const data = JSON.parse(event.data);
      for (let i = 0; i < PLAYER_NUM-1; i++) {
        if (data[i] == null) continue;

        const info = data[i];

        const id = `enemy-${info['player_id']}`;
        const enemy = this.enemies.find((e) => e.id == id);

        let diff = 0;
        let sTime = sTimes[i];
        if (sTime == -1) {
          sTime = Date.now();
        } else {
          const sTime1 = Date.now();
          diff = sTime1 - sTime;
          sTime = sTime1;
        }
        sTimes[i] = sTime;

        enemy.registerAction(0, Number(info['direction']), Number(info['position']['x']), Number(info['position']['y']));

        enemy.interpolateLag2(diff, Number(info['direction']), Number(info['position']['x']), Number(info['position']['y']))

        enemy.setDirectionPhi(Number(info['direction']));
        enemy.setPos(Number(info['position']['x']), Number(info['position']['y']));
      }
    };
  }

  // TODO: 前以外の移動に対応する

  createEnemies(n: number) {
    this.enemies = Array<Enemy>(n);
    for (let i = 0; i < n; i++) {
      const enemy = new Enemy(this, `enemy-${i+1}`, 0, 0);
      this.add.existing(enemy);
      this.physics.world.enable(enemy);
      this.enemies[i] = enemy;
    }
    this.player.setOverlap(this.enemies);
    for (const enemy of this.enemies) {
      const other: Array<Ship> = this.enemies.filter((e) => e.id != enemy.id);
      other.push(this.player as Ship);
      enemy.setOverlap(other);
    }
  }

  createMap(width: number, height: number) {
    const graphics = this.add.graphics();

    // draw Outline
    graphics.lineStyle(5, 0x00fd00, 1);
    graphics.beginPath();
    graphics.moveTo(0, 0);
    graphics.lineTo(0, height);
    graphics.lineTo(width, height);
    graphics.lineTo(width, 0);
    graphics.closePath();
    graphics.strokePath();

    // draw Mesh
    graphics.lineStyle(1, 0x00fd00, 1);
    graphics.beginPath();

    for (let i = 1; i <= 10; i++) {
      graphics.moveTo(width * i / 10, 0);
      graphics.lineTo(width * i / 10, height);
    }
    for (let i = 1; i <= 10; i++) {
      graphics.moveTo(0, height * i / 10);
      graphics.lineTo(width, height * i / 10);
    }
    graphics.strokePath();
  }

  update() {
    const alives = this.enemies.filter((e) => e.active);
    if (!this.player.active || alives.length == 0) {
      this.scene.stop('game');
      this.scene.run('game-over', {ranking: alives.length + 1});
    }

    text.setText([
      'Ranking: ' + (alives.length + 1),
    ]);

    for (const enemy of this.enemies) {
      enemy.interpolateLag();
    }
  }
}
