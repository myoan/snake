import 'phaser';
import Game from './src/scene/Game';
import Preloader from './src/scene/Preloader';


const config = {
  type: Phaser.AUTO,
  width: 1200,
  height: 1000,
  parent: 'phaser-example',
  physics: {
    default: 'arcade',
  },
  scene: [ Preloader, Game ],
};

new Phaser.Game(config);
