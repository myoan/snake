import axios from 'axios';
import 'phaser';

let text: Phaser.GameObjects.Text;
export default class Preloader extends Phaser.Scene {
  id: String;
  graphics: Phaser.GameObjects.Graphics;
  conn: WebSocket;
  ip: string;
  port: integer;

  constructor(id: String, score: integer) {
    super('preloader');
  }

  create(args) {
    const id = args[0];
    const score = args[1] | 0;
    const content = [
      `ID: ${id}`,
      `Score: ${score}`,
      "[ENTER] -> Game Start"
    ]
    this.id = id;

    text = this.add.text(100, 100, content, { fontFamily: 'Arial', color: '#00ff00' });
    this.input.keyboard.on('keydown-ENTER', () => { this.connect() }, this);
  }

  connect() {
    const scene = this;
    console.log('get api.snake.game.myoan.dev/room');

    axios.get('https://api.snake.game.myoan.dev/room')
      .then(function (resp) {
        const data = resp.data;

        scene.ip = data.ip;
        scene.port = data.port;

        console.log("connect " + scene.ip + ':' + scene.port);
        scene.conn = new WebSocket('wss://ws.snake.game.myoan.dev:' + scene.port + '/');
        // scene.conn = new WebSocket('wss://' + scene.ip + ":" + scene.port + '/');
        // scene.conn = new WebSocket('wss://' + scene.ip + ":" + scene.port + '/');

        scene.conn.onmessage = (event) => {
          const data = JSON.parse(event.data);
          switch(data.status) {
            case 0: // GameStatusInit
              scene.id = data.id
              break;

            case 1: // GameStatusOk
              scene.scene.start('game', [scene.id, scene.conn])
              break

            case 3: // GameStatusWaiting...
              console.log(`waiting ...`)

              text.destroy();
              text = scene.add.text(100, 100, "waiting...", { fontFamily: 'Arial', color: '#00ff00' });
              break;

            default:
              console.log(`data(${data.status}): ${data}`)
          }
        };
      })
      .catch(function (error) {
        console.log(error);
      });
  }
}
