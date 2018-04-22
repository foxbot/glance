// -- Global
const ApiUrl = "ws://" + window.location.host + "/api/socket";
let TotalShards = 0;

/** @type {Array<HTMLDivElement>} */
let shards = [];

// -- Shard Calculator

const RightShift = 4194304;

/** @type {HTMLInputElement} */
const shardInput = document.getElementById("shard-calc-in");
/** @type {HTMLSpanElement} */
const shardOutput = document.getElementById("shard-calc-out");

shardInput.addEventListener("input", function(event) {
  const id = +event.target.value;
  const shardNumber = Math.floor(id / RightShift) % TotalShards;
  shardOutput.innerText = shardNumber;
  setSelectedShard(shardNumber);
});

function setSelectedShard(id) {
  shards.forEach(shard => {
    if (shard.classList.contains("selected")) {
      shard.classList.remove("selected");
    }
  });
  shards[id].classList.add("selected");
}

// -- Shards Container

const shardsContainer = document.getElementById("shards-container");

function paintShards() {
  for (let i = 0; i < TotalShards; i++) {
    const el = document.createElement("div");
    el.className = "shard";
    el.setAttribute("data-status", "unknown");
    
    const nameEl = document.createElement("h3");
    nameEl.className = "shard-id";
    nameEl.innerText = i;
  
    el.appendChild(nameEl);
  
    shardsContainer.appendChild(el);
    shards[i] = el;
  }
}
paintShards();

// -- Socket
let socket = new WebSocket(ApiUrl);
socket.addEventListener("open", function(e) {
  console.log("socket open", e);
});
socket.addEventListener("close", function(e) {
  console.log("socket close", e);
});
socket.addEventListener("error", function(e) {
  console.log("socket error", e);
});

const OP = {
  HELLO: 0
};

socket.addEventListener("message", function(e) {
  console.log("socket message", e);
  
  let data = JSON.parse(e.data);
  if (!'Op' in data) {
    console.error("invalid payload", data);
    return;
  }
  switch (data.Op) {
    case OP.HELLO:
      TotalShards = data.Data.TotalShards;
      paintShards();
      break;
  }
});
