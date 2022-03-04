export class Vector {
  constructor(public x: number, public y: number) {
  }

  magnitude(): number {
    return Math.sqrt(this.x * this.x + this.y * this.y);
  }

  normalize(): Vector {
    const x = this.x / this.magnitude();
    const y = this.y / this.magnitude();
    return new Vector(x, y);
  }

  add(v:Vector) {
    const x = this.x + v.x;
    const y = this.y + v.y;
    return new Vector(x, y);
  }

  sub(v:Vector): Vector {
    const x = this.x - v.x;
    const y = this.y - v.y;
    return new Vector(x, y);
  }

  mul(n: number): Vector {
    const x = this.x * n;
    const y = this.y * n;
    return new Vector(x, y);
  }

  div(n: number): Vector {
    const x = this.x / n;
    const y = this.y / n;
    return new Vector(x, y);
  }

  show(): string {
    return `(${this.x}, ${this.y})`;
  }

  cosign(v: Vector): number {
    if (this.x == 0 && this.y == 0) throw new Error('object is not vector');
    if (v.x == 0 && v.y == 0) throw new Error('args is not vector');

    return this.innerProduct(v) / (Math.sqrt(this.x * this.x + this.y * this.y) * Math.sqrt(v.x * v.x + v.y * v.y));
  }

  innerProduct(v: Vector): number {
    return this.x * v.x + this.y * v.y;
  }

  angle(v: Vector): number {
    if (this.x == 0 && this.y == 0) throw new Error('object is not vector');
    if (v.x == 0 && v.y == 0) throw new Error('args is not vector');

    const cos = this.normalize().cosign(v.normalize());
    return Math.acos(cos) * 180 / Math.PI;
  }
}

export default Vector;
