import 'phaser';

const WINDOW_WIDTH = 1200;
const WINDOW_HEIGHT = 1000;
const CELL_PX = 16;

export class Board {
  scene: Phaser.Scene;
  width: number;
  height: number;
  raw: integer[][];

  constructor(scene: Phaser.Scene, w: number, h: number) {
    this.scene = scene;
    this.width = w;
    this.height = h;

    this.raw = new Array();
    for (let i = 0; i < this.height; i++) {
      this.raw[i] = new Array(this.width).fill(0);
    }
  }

  setCell(x: number, y: number, v: number) {
    this.raw[y][x] = v;
  }

  sync(data: integer[][]) {
    for(let i = 0; i < this.height; i++) {
      for (let j = 0; j < this.width; j++) {
        this.raw[i][j] = data[i][j];
      }
    }
  }

  draw(data: integer[][], force: Boolean = false) {
    if (!force && data.length == 0) {
      data = this.raw
    }

    const widthPx = this.width * (CELL_PX+4)
    const heightPx = this.height * (CELL_PX+4)

    const xPad = (WINDOW_WIDTH - widthPx) / 2
    const yPad = (WINDOW_HEIGHT - heightPx) / 2

    for (let i = 0; i < this.height; i++) {
      for (let j = 0; j < this.width; j++) {
        if (!force && this.raw[i][j] == data[i][j]) {
          continue;
        }
        this.raw[i][j] = data[i][j]
        const x = xPad + j * (CELL_PX+4);
        const y = yPad + i * (CELL_PX+4);
        if (this.raw[i][j] < 0) {
          this.scene.add.rectangle(x, y, CELL_PX, CELL_PX, 0xff9999);
        } else if (this.raw[i][j] > 0) {
          this.scene.add.rectangle(x, y, CELL_PX, CELL_PX, 0xcccccc);
        } else {
          this.scene.add.rectangle(x, y, CELL_PX, CELL_PX, 0x000000);
          var cell = this.scene.add.rectangle(x, y, CELL_PX, CELL_PX);
          cell.setStrokeStyle(1, 0x1a65ac);
        }
      }
    }
  }
}

export default Board;
