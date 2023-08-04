import Phaser from "phaser";

const icons = [
  "shinomiya",
  "kaicho",
  "hayasaka",
  "keichang",
];
const suits = ["s", "h", "d", "c"];
const nums = [
  "-",
  "A",
  "2",
  "3",
  "4",
  "5",
  "6",
  "7",
  "8",
  "9",
  "T",
  "J",
  "Q",
  "K",
];
function num(code) {
  console.assert(code !== -1);
  return nums[code & 0b1111];
}
function suit(code) {
  console.assert(code !== -1);
  return suits[(code >> 4) & 0b11];
}
function revealed(code) {
  console.assert(code !== -1);
  return (code >> 6) & 1;
}
function reveal(code) {
  console.assert(code !== -1);
  return code | (1 << 6);
}

export class Main extends Phaser.Scene {
  constructor() {
    super("Main");
  }
  preload() {
    this.load.plugin(
      "rexperspectiveimageplugin",
      "https://raw.githubusercontent.com/rexrainbow/phaser3-rex-notes/master/dist/rexperspectiveimageplugin.min.js",
      true,
    );
    this.load.plugin(
      "rexcirclemaskimageplugin",
      "https://raw.githubusercontent.com/rexrainbow/phaser3-rex-notes/master/dist/rexcirclemaskimageplugin.min.js",
      true,
    );
    icons.forEach((icon) => {
      this.load.image(icon, `icons/${icon}.jpg`);
    });
    suits.forEach((suit) =>
      nums.slice(1).forEach((num) => {
        const card = suit + num;
        this.load.image(card, `cards/${card}.png`);
      })
    );
    this.load.image("back", "cards/back.png");
    this.load.image("empty", "cards/empty.png");
    this.canvas = this.sys.game.canvas;
    this.cfg = {
      icon: [
        {
          x: this.canvas.width / 2 - 220,
          y: this.canvas.height / 2 + 130,
        },
        {
          x: this.canvas.width / 2 - 420,
          y: this.canvas.height / 2,
        },
        {
          x: this.canvas.width / 2 + 200,
          y: this.canvas.height / 2 - 140,
        },
        {
          x: this.canvas.width / 2 + 420,
          y: this.canvas.height / 2,
        },
      ],
      deck: {
        x: this.canvas.width / 2 + 100,
        y: this.canvas.height / 2,
      },
      trash: {
        x: this.canvas.width / 2 - 100,
        y: this.canvas.height / 2,
      },
      // button: {
      //   fold: {
      //     x: this.canvas.width - 250,
      //     y: this.canvas.height - 50,
      //   },
      //   counter: {
      //     x: this.canvas.width - 100,
      //     y: this.canvas.height - 50,
      //   },
      // },
      hands: [
        {
          x: this.canvas.width / 2,
          y: this.canvas.height - 100,
          align: {
            dx: 1.7,
            dy: 0,
          },
          msg: {
            dx: 0,
            dy: -1,
          },
          angle: 0,
        },
        {
          x: 100,
          y: this.canvas.height / 2,
          align: {
            dx: 0,
            dy: 1,
          },
          msg: {
            dx: 2.2,
            dy: 0,
          },
          angle: -90,
        },
        {
          x: this.canvas.width / 2,
          y: 100,
          align: {
            dx: -1,
            dy: 0,
          },
          msg: {
            dx: 0,
            dy: 1,
          },
          angle: 0,
        },
        {
          x: this.canvas.width - 100,
          y: this.canvas.height / 2,
          align: {
            dx: 0,
            dy: -1,
          },
          msg: {
            dx: -2.2,
            dy: 0,
          },
          angle: 90,
        },
      ],
    };
  }
  clear() {
    this.hands.forEach((hand) => {
      hand.forEach((card) => card.destroy());
    });
    this.hands = [[], [], [], []];
    this.trash.forEach((card) => card.destroy());
    this.trash = [];
    this.moriQueue = [];
    this.clearChat();
  }
  clearChat() {
    this.chatbox.forEach((chat) => {
      if (chat !== null) chat.destroy();
    });
    this.chatbox = [null, null, null, null];
  }
  createChat(player, msg) {
    const cfg = this.cfg.hands[player];
    const chat = this.add.text(
      cfg.x + cfg.msg.dx * 130,
      cfg.y + cfg.msg.dy * 130,
      msg,
      {
        fontSize: "40px",
      },
    ).setOrigin(0.5)
      .setTint(0x020202);

    const old = this.chatbox[player];
    if (old !== null) {
      old.destroy();
    }
    this.chatbox[player] = chat;
  }

  Fetch(data) {
    this.clear();
    // this.createChat(0, "MORI!!");
    // this.createChat(1, "MORI!!");
    // this.createChat(2, "MORI!!");
    // this.createChat(3, "MORI!!");

    // if (data.top !== -1) {
    const top = this.createCard({
      x: this.cfg.trash.x,
      y: this.cfg.trash.y,
      name: (data.top !== -1 ? suit(data.top) + num(data.top) : "empty"),
      face: 0,
      angle: 0,
    });
    this.trash.push(top);
    // }

    data.hands[this.id].forEach((code) => {
      const card = this.createMyCard(code);
      this.hands[0].push(card);
    });
    this.arange(0);

    for (let i = 1; i <= 3; i++) {
      const p = (this.id + i) % 4;
      data.hands[p].forEach((code) => {
        const card = this.createOthersCard(code);
        this.hands[i].push(card);
      });
      this.arange(i);
    }

    data.moriQueue.forEach((p) => {
      const i = (p - this.id + 4) % 4;
      this.createChat(i, "MORI!!");
      this.moriQueue.push(i);
    });
  }

  Discard(data) {
    const i = (data.player - this.id + 4) % 4;
    var card = this.hands[i].find((card) => card.code === data.card);

    if (card === undefined) {
      const idx = this.hands[i].findIndex((card) => card.code === -1);
      const discard = this.hands[i][idx];
      card = this.createCard({
        x: discard.x,
        y: discard.y,
        name: suit(data.card) + num(data.card),
        face: 1,
        angle: discard.angle,
      });
      card.setInteractive().on("pointerup", () => {
        this.ws.send(`{"type":"mori"}`);
      });
      discard.destroy();
      this.hands[i][idx] = card;
    }

    if (i === 0) {
      card.setInteractive().on("pointerup", () => {
        this.ws.send(`{"type":"counter"}`);
      });
    }

    this.moveToTrash(card);
    card.setAlpha(1.0);

    this.hands[i] = this.hands[i].filter((h) => h !== card);
    this.arangeMove(i);
  }

  drawMyHand(code) {
    const draw = this.createMyCard(code);
    draw.setFace(1);
    draw.flip.flip();
    this.hands[0].push(draw);
    this.arangeMove(0);
  }
  drawOthersHand(player, code) {
    const draw = this.createOthersCard(code);
    this.hands[player].push(draw);
    console.log("player:", player);
    console.log("this.hands[player]:", this.hands[player]);
    this.arangeMove(player);
  }
  Draw(data) {
    this.moriQueue.forEach((i) => {
      this.hands[i].forEach((card) => {
        this.moveToTrash(card);
        card.setAlpha(1.0);
      });
      this.hands[i] = [];
    });
    this.moriQueue = [];

    const i = (data.player - this.id + 4) % 4;
    if (i === 0) {
      this.drawMyHand(data.card);
    } else {
      this.drawOthersHand(i, data.card);
    }
  }

  Flip(data) {
    this.clearChat();
    const top = this.createCard({
      ...this.cfg.deck,
      name: suit(data.card) + num(data.card),
      face: 1,
    });
    this.moveToTrash(top);
    const i = (data.player - this.id + 4) % 4;
    if (i !== 0) {
      top.setInteractive().on("pointerup", () => {
        this.ws.send(`{"type":"mori"}`);
      });
    } else {
      top.setInteractive().on("pointerup", () => {
        this.ws.send(`{"type":"counter"}`);
      });
    }
  }

  Mori(data) {
    const i = (data.player - this.id + 4) % 4;
    this.createChat(i, "MORI!!");
    this.moriQueue.push(i);
    if (i !== 0) {
      this.revealOthersHand(i, data.cards);
    }
  }

  Fold(data) {
    const i = (data.player - this.id + 4) % 4;
    this.createChat(i, "FOLD..");
  }

  Counter(data) {
    const i = (data.player - this.id + 4) % 4;
    this.createChat(i, "wwwwww");
    if (i !== 0) {
      this.revealOthersHand(i, data.cards);
    }
    const loser = (data.loser - this.id + 4) % 4;
    this.createChat(loser, "F**K");
    this.moriQueue.push(i);
  }

  revealMyHand() {
    this.hands[0].forEach((card) => {
      card.setAlpha(0.7);
      card.code = reveal(card.code);
      card.setInteractive().on("pointerup", () => {
        this.ws.send(`{"type":"discard","card":${card.code}}`);
      });
    });
  }
  revealOthersHand(player, cards) {
    cards.forEach((code) => {
      const idx = this.hands[player].findIndex((c) => c.code === -1);
      const card = this.hands[player][idx];
      const newCard = this.createCard({
        x: card.x,
        y: card.y,
        name: suit(code) + num(code),
        face: 1,
        angle: card.angle,
      });
      card.destroy();
      newCard.code = reveal(code);
      newCard.flip.flip();
      this.hands[player][idx] = newCard;
    });
    this.arangeMove(player);
  }
  Reveal(data) {
    const i = (data.player - this.id + 4) % 4;
    if (i === 0) {
      this.revealMyHand();
    } else {
      this.revealOthersHand(i, data.cards);
    }
  }
  Burst(data) {
    const i = (data.player - this.id + 4) % 4;
    if (i !== 0) {
      this.revealOthersHand(i, [data.card]);
    }

    this.hands[i].forEach((card) => {
      this.tweens.add({
        targets: card,
        x: this.cfg.trash.x,
        y: this.cfg.trash.y,
        ease: "Cubic",
        angle: card.angle,
      });
      this.trash.push(card);
    });
    this.hands[i] = [];
  }

  setupIcons() {
    icons.forEach((icon, i) => {
      const a = this.make.rexCircleMaskImage({
        x: this.cfg.icon[(i - this.id + 4) % 4].x,
        y: this.cfg.icon[(i - this.id + 4) % 4].y,
        key: icon,
        maskType: 0,
      });
      a.setScale(0.8);
    });
  }

  create() {
    const searchParams = new URLSearchParams(window.location.search);
    const hostname = window.location.hostname;
    const protocol = window.location.protocol;
    this.ws = new WebSocket(
      `${protocol === "http:" ? "ws" : "wss"}://${hostname}/ws?room=${
        searchParams.get("room")
      }`,
    );
    this.ws.onmessage = (e) => {
      const data = JSON.parse(e.data);
      console.log("data:", data);
      switch (data.type) {
        case "join":
          this.id = data.id;
          this.setupIcons();
          break;
        case "fetch":
          this.Fetch(data);
          break;
        case "discard":
          this.Discard(data);
          break;
        case "draw":
          this.Draw(data);
          break;
        case "flip":
          this.Flip(data);
          break;
        case "mori":
          this.Mori(data);
          break;
        case "fold":
          this.Fold(data);
          break;
        case "counter":
          this.Counter(data);
          break;
        case "reveal":
          this.Reveal(data);
          break;
        case "burst":
          this.Burst(data);
          break;
        case "shuffle":
          this.Shuffle(data);
          break;
        default:
          console.log("Invalid type");
      }
    };

    this.hands = [[], [], [], []];
    this.moriQueue = [];
    this.chatbox = [null, null, null, null];
    const deck = this.createCard({
      ...this.cfg.deck,
      face: 1,
    });
    deck.setInteractive().on("pointerup", () => {
      this.ws.send(`{"type":"fold"}`);
      this.ws.send(`{"type":"flip"}`);
      this.ws.send(`{"type":"draw"}`);
    });
    this.trash = [];
    const top = this.createCard({
      x: this.cfg.trash.x,
      y: this.cfg.trash.y,
      name: "empty",
      face: 0,
      angle: 0,
    });
    this.trash.push(top);
    // this.buttons = [];
    // this.createButton(
    //   this.cfg.button.fold.x,
    //   this.cfg.button.fold.y,
    //   "FOLD",
    //   () => {
    //     this.ws.send(`{"type":"fold"}`);
    //   },
    // );
    // this.createButton(
    //   this.cfg.button.counter.x,
    //   this.cfg.button.counter.y,
    //   "COUNTER",
    //   () => {
    //     this.ws.send(`{"type":"counter"}`);
    //   },
    // );
  }
  createButton(x, y, text, onClick) {
    const button = this.add.text(
      x,
      y,
      text,
      {
        fontSize: "30px",
      },
    ).setOrigin(0.5)
      .setStyle({ backgroundColor: "#333" })
      .setPadding(10)
      .setInteractive({ useHandCursor: true })
      .on("pointerup", onClick);
    this.buttons.push(button);
  }
  createMyCard(code) {
    const card = this.createCard({
      ...this.cfg.deck,
      name: suit(code) + num(code),
      transparent: revealed(code),
    });
    card.code = code;
    card.setInteractive().on("pointerup", () => {
      this.ws.send(`{"type":"discard","card":${code}}`);
    });
    return card;
  }
  createOthersCard(code) {
    console.log("createOthersCard");
    const card = this.createCard({
      ...this.cfg.deck,
      name: code === -1 ? "empty" : suit(code) + num(code),
      face: code === -1 ? 1 : 0,
    });
    card.code = code;
    card.setInteractive().on("pointerup", () => {
      this.ws.send(`{"type":"mori"}`);
    });
    return card;
  }
  moveToTrash(card) {
    this.tweens.add({
      targets: card,
      x: this.cfg.trash.x + (Math.random() - 0.5) * 50,
      y: this.cfg.trash.y + (Math.random() - 0.5) * 50,
      ease: "Cubic",
      angle: card.angle + (Math.random() - 0.5) * 20,
    });
    this.trash.push(card);
    if (card.face === 1) card.flip.flip();
    card.setDepth(10000 + this.trash.length);
  }
  createCard(cfg) {
    let card = this.add.rexPerspectiveCard({
      x: cfg.x,
      y: cfg.y,
      front: { key: cfg.name || "empty" },
      back: { key: "back" },
      face: cfg.face || 0,
    });
    card.setScale(cfg.scale || 3.0);
    card.setAngle(cfg.angle || 0);
    card.setAlpha(cfg.transparent ? 0.7 : 1.0);
    return card;
  }
  arangeMove(player) {
    console.log("arangeMove");
    const cfg = this.cfg.hands[player];
    const hand = this.hands[player];
    hand.forEach((h, idx) => {
      this.tweens.add({
        targets: h,
        x: cfg.x + cfg.align.dx * 80 * (idx - (hand.length - 1) / 2),
        y: cfg.y + cfg.align.dy * 80 * (idx - (hand.length - 1) / 2),
        angle: cfg.angle,
        ease: "Cubic",
      });
    });
  }
  arange(player) {
    const cfg = this.cfg.hands[player];
    const hand = this.hands[player];
    hand.forEach((h, idx) => {
      h.setX(cfg.x + cfg.align.dx * 80 * (idx - (hand.length - 1) / 2));
      h.setY(cfg.y + cfg.align.dy * 80 * (idx - (hand.length - 1) / 2));
      h.setAngle(cfg.angle);
    });
  }
}
