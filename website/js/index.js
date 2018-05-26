// -- Global
const ApiUrl = (window.location.protocol == "https:" ? "wss://" : "ws://") + window.location.host + "/api/socket";
let TotalShards = 0;
let TotalPatrons = 0;
let PatronShards = 0;


let shards = {
  /** @type {Array<HTMLDivElement>} 
   * Main bot shards
  */
  0: [],
  /** @type {Array<HTMLDivElement>}
   * Patron Bot 1
   */
  1: [],
  /** @type {Array<HTMLDivElement>} 
   * Patron Bot 2
  */
  2: [],
};

// -- Shard Calculator

const RightShift = 4194304;

/** @type {HTMLInputElement} */
const shardInput = document.getElementById("shard-calc-in");
/** @type {HTMLSpanElement} */
const shardOutput = document.getElementById("shard-calc-out");
/** @type {HTMLSpanElement} */
const shardOutputPatron = document.getElementById("shard-calc-pout");

shardInput.addEventListener("input", function(event) {
  const id = +event.target.value;
  if (isNaN(id)) {
    setSelectedShard(-1, -1);
    return;
  }

  const shardNumber = Math.floor(id / RightShift) % TotalShards;
  shardOutput.innerText = shardNumber;
  const shardNumberP = Math.floor(id / RightShift) % PatronShards;
  shardOutputPatron.innerText = shardNumberP;
  setSelectedShard(shardNumber, shardNumberP);
});

function setSelectedShard(id, patronId) {
  for (let bot in shards) {
    shards[bot].forEach(shard => {
      if (shard.classList.contains("selected")) {
        shard.classList.remove("selected");
      }
    });
    if (bot == 0 && id > 0) {
      shards[0][id].classList.add("selected");
    }
    if (bot > 0 && patronId > 0) {
      shards[+bot][patronId].classList.add("selected");
    }
  }
}

// -- Shards Container

const shardsContainer = document.getElementById("shards-container");
const patronContainer = document.getElementById("patron-container");

function createShardEl(name) {
  const el = document.createElement("div");
  el.className = "shard";
  el.setAttribute("data-status", "unknown");
    
  const nameEl = document.createElement("p");
  nameEl.className = "shard-id";
  nameEl.innerText = name;
  
  el.appendChild(nameEl);

  return el;
}

function paintShards() {
  while (shardsContainer.lastChild) {
    shardsContainer.removeChild(shardsContainer.lastChild)
  }
  while (patronContainer.lastChild) {
    patronContainer.removeChild(patronContainer.lastChild);
  }
  
  for (let i = 0; i < TotalShards; i++) {
    let el = createShardEl(i);
  
    shardsContainer.appendChild(el);
    shards[0][i] = el;
  }
  for (let p = 1; p < TotalPatrons+1; p++) {
    for (let i = 0; i < PatronShards; i++) {
      let el = createShardEl("P"+p+"-"+i);
      
      patronContainer.appendChild(el);
      shards[p][i] = el;
    }
  }
}
paintShards();

// -- Socket

let socket = new WebSocket(ApiUrl);
addEvents(socket);

function addEvents(socket) {
  socket.addEventListener("open", onOpen);
  socket.addEventListener("close", onClose);
  socket.addEventListener("error", onError);
  socket.addEventListener("message", onMessage);
}

function onOpen(e) {
  console.log("Socket opened.");
}
function onClose(e) {
  console.log("Socket closed.");
  for (let key in shards) {
    shards[key] = [];
  }
  TotalShards = 0;
  paintShards();
  
  // Try to reconnect
  // TODO: exponential backoff? probably should.
  setTimeout(function() {
    console.log("Attempting reconnect...");
    socket = new WebSocket(ApiUrl);
    addEvents(socket);
  }, 5000);
}
function onError(e) {
  console.error("Socket error!", e);
}

const OP = {
  HELLO: 0,
  UPDATE: 1,
  TICK: 2,
};
const STATUS = {
  0: 'unknown',
  1: 'waiting',
  2: 'starting',
  3: 'online',
  4: 'stopping',
  5: 'offline'
};

function onMessage(e) {
  let data = JSON.parse(e.data);
  if (!'Op' in data) {
    console.error("invalid payload", data);
    return;
  }
  switch (data.Op) {
    case OP.HELLO:
      TotalShards = data.Data.TotalShards;
      TotalPatrons = data.Data.TotalPatrons;
      PatronShards = data.Data.PatronShards;
      
      paintShards();

      for (const id in data.Data.State[0] ) {
        const val = STATUS[data.Data.State[0][id]];
        shards[0][id].setAttribute("data-status", val);
      }

      break;
    case OP.UPDATE:
      const bot = +data.Data.Bot;
      const id = +data.Data.Shard;
      const val = STATUS[data.Data.Status];
      shards[bot][id].setAttribute("data-status", val)
      break;
    // don't care about heartbeat data, just keep the socket alive
    case OP.TICK:
      break;
  }
}
