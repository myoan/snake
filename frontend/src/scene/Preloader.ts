import 'phaser';
import axios from 'axios';

let text: Phaser.GameObjects.Text;
export default class Preloader extends Phaser.Scene {
  id: String;
  graphics: Phaser.GameObjects.Graphics;
  conn: WebSocket;
  dotNum: integer;
  tick: integer;

  constructor(id: String, score: integer) {
    super('preloader');
  }

  create(args) {
    const id = args[0];
    const score = args[1] | 0;
    const content = [
      `ID: ${id}`,
      `Score: ${score}`,
      "[ENTER] -> Game Start!!"
    ]
    this.id = id;

    text = this.add.text(100, 100, content, { fontFamily: 'Arial', color: '#00ff00' });
    console.log('get api.snake-game.myoan.dev/room');

    const scene = this;

    axios.get('http://api.snake-game.myoan.dev/room')
      .then(function (resp) {
        console.log(resp);
        const data = resp.data;

        scene.input.keyboard.on('keydown-ENTER', () => { scene.connect(data.ip, data.port) }, scene);
      })
      .catch(function (error) {
        console.log(error);
      });
  }

  connect(ip: String, port: integer) {
    console.log("connect " + ip + port.toString());
    this.conn = new WebSocket('ws://' + ip + ":" + port.toString() + '/');

    this.conn.onmessage = (event) => {
      const data = JSON.parse(event.data);
      switch(data.status) {
        case 0: // GameStatusInit
          this.id = data.id
          break;

        case 1: // GameStatusOk
          this.scene.start('game', [this.id, this.conn])
          break

        case 3: // GameStatusWaiting...
          console.log(`waiting ...`)

          text.destroy();
          text = this.add.text(100, 100, "waiting...", { fontFamily: 'Arial', color: '#00ff00' });
          break;

        default:
          console.log(`data(${data.status}): ${data}`)
      }
    };
  }
}
