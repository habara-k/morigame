import Phaser from "phaser";

import { Main } from "./scenes/Main";

const config = {
  width: 1334,
  height: 750,
  type: Phaser.AUTO,
  backgroundColor: 0xcdcdcd,

  scale: {
    mode: Phaser.Scale.FIT,
    autoCenter: Phaser.Scale.CENTER_BOTH,
  },

  scene: Main,
  antialias: false,
};

const game = new Phaser.Game(config);
export default game;
