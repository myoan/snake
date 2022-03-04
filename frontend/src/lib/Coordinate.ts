import Vector from './Vector';

export default class Coordinate {
  pos: Vector;
  theta: number;
  constructor(pos: Vector, theta: number) {
    this.pos = pos;
    this.theta = theta;
  }

  move(r: number, phai: number = 0) {
    const rad = (this.theta + phai) * (Math.PI / 180);
    const newX = this.pos.x + r * Math.cos(rad);
    const newY = this.pos.y + r * Math.sin(rad);
    this.pos = new Vector(newX, newY);
  }

  rotate(d: number) {
    this.theta = (this.theta + d) % 360;
  }

  convertToWorld(pos: Vector): Vector {
    const rad = this.theta * (Math.PI / 180);
    const x = this.pos.x + Math.cos(rad) * pos.x - Math.sin(rad) * pos.y;
    const y = this.pos.y + Math.sin(rad) * pos.x + Math.cos(rad) * pos.y;
    return new Vector(x, y);
  }

  convertToLocal(pos: Vector): Vector {
    const rad = this.theta * (Math.PI / 180);
    const sin = Math.sin(rad);
    const cos = Math.cos(rad);
    const denom = sin * sin + cos * cos;
    const x = -1 * ((this.pos.x - pos.x) * cos + (this.pos.y - pos.y) * sin) / denom;
    const y = ((this.pos.x - pos.x) * sin - (this.pos.y - pos.y) * cos) / denom;
    return new Vector(x, y);
  }

  directionToWorld(phai: number, r: number = 1): Vector {
    const rad = (this.theta + phai) * (Math.PI / 180);
    const x = r * Math.cos(rad);
    const y = r * Math.sin(rad);
    return new Vector(x, y);
  }
}
