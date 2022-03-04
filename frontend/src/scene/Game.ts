import Board from '../game/Board';

const SCREEN_WIDTH_CELL = 40;
const SCREEN_HEIGHT_CELL = 40;

const MOVE_LEFT = 0;
const MOVE_RIGHT = 1;
const MOVE_UP = 2;
const MOVE_DOWN = 3;

let text: Phaser.GameObjects.Text;

function arrayTo2DArray(arry: integer[], width: integer, height: integer): integer[][] {
  const result = new Array();
  for (let i = 0; i < height; i++) {
    result[i] = new Array(width).fill(0)
  }
  for (let i = 0; i < height; i++) {
    for (let j = 0; j < width; j++) {
      result[i][j] = arry[i*height + j]
    }
  }
  return result
}
export default class Game extends Phaser.Scene {
  id: String;
  board: Board;
  conn: WebSocket;
  constructor(conn: WebSocket) {
    super('game')
  }

  create(args) {
    var id = args[0];
    var conn = args[1];
    this.id = id;
    this.conn = conn;
    this.board = new Board(this, SCREEN_WIDTH_CELL, SCREEN_HEIGHT_CELL);
    const initArray = new Array(SCREEN_HEIGHT_CELL * SCREEN_WIDTH_CELL).fill(0);
    this.board.draw(arrayTo2DArray(initArray, SCREEN_WIDTH_CELL, SCREEN_HEIGHT_CELL), true);

    this.input.keyboard.on('keydown-W', () => { this.sendDirection(MOVE_UP) }, this)
    this.input.keyboard.on('keydown-A', () => { this.sendDirection(MOVE_LEFT) }, this)
    this.input.keyboard.on('keydown-S', () => { this.sendDirection(MOVE_DOWN) }, this)
    this.input.keyboard.on('keydown-D', () => { this.sendDirection(MOVE_RIGHT) }, this)

    this.conn.onmessage = (event) => {
      const data = JSON.parse(event.data);
      switch(data.status) {
        case 1: // GameStatusOk
          const body = data.body
          this.board.draw(arrayTo2DArray(body.board, SCREEN_WIDTH_CELL, SCREEN_HEIGHT_CELL));
          break;

        case 2: // GameStatusError
          console.log(`dropped`)
          this.conn.close();
          data.body.players.forEach(p => {
            if (p.id == this.id) {
              this.scene.start('preloader', [this.id, p.size])
            }
          })
          break;

        default:
          console.log(`data: ${data}`)
      }
    };
  }

  sendDirection(dir: integer) {
    const data =  {
      uuid: this.id,
      eventtype: 0,
      key: dir
    };
    var json = JSON.stringify(data);
    // console.log(`send data: ${json}`);
    this.conn.send(json);
  }
}
