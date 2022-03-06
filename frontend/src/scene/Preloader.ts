import 'phaser';

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
      "[ENTER] -> Game Start"
    ]
    this.id = id;
    text = this.add.text(100, 100, content, { fontFamily: 'Arial', color: '#00ff00' });

    this.input.keyboard.on('keydown-ENTER', () => { this.connect() }, this);
  }

  connect() {
    console.log("connect")
    this.conn = new WebSocket('ws://192.168.49.2:31290/ingame');
    // this.conn = new WebSocket('ws://' + document.location.hostname + ':31290/ingame');

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
