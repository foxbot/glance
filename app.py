from enum import IntEnum

from dotenv import load_dotenv
from sanic import Sanic
from sanic.response import json, file, text
from sanic.websocket import WebSocketProtocol

import asyncio
import logging
import os
import time
import ujson


class ShardState(IntEnum):
    UNKNOWN = 0
    UNHEALTHY = 1
    OFFLINE = 2
    ONLINE = 3
    STARTING = 4
    STOPPING = 5
    RESUMING = 6
    # TODO: are these good states, needs to check with discord.py

# {
#   "bot": {
#     "len": 1,
#     "0": {
#       "state": 0
#       "last_healthy": 0
#     }
#   }
# }
State = dict()
Feeds = list()
app = Sanic()
HealthTask = None
health_timeout = None

app.static('/assets/', './assets/')


def validate_shard_in_cluster(cluster: str, id: int):
    if cluster not in State:
        return json({
            'error': f'cluster \'{cluster}\' does not exist'
        }, status=400)
    
    if id < 0 or id > State[cluster]['len']-1:
        return json({
            'error': f'shard \'{id}\' falls outside the range of valid shards on cluster \'{cluster}\''
        }, status=400)
    
    return None


@app.route('/api/health/<cluster>/<id:int>', methods=['POST'])
async def health(request, cluster: str, id: int):
    valid = validate_shard_in_cluster(cluster, id)
    if valid is not None: return valid
    
    if str(id) not in State[cluster]:
        return json({
            'error': f'shard \'{id}\' is not known to the system yet, please populate your status'
        }, status=409)
    
    State[cluster][str(id)]['last_healthy'] = time.time()
    msg = {
        'op': 'health_ping',
        'payload': {
            'cluster': cluster,
            'shard': id
        }
    }
    
    if State[cluster][str(id)]['state'] == ShardState.UNHEALTHY:
        # assume shards won't health-check unless they are online :)
        State[cluster][str(id)]['state'] = ShardState.ONLINE
        msg['payload']['state'] = ShardState.ONLINE

    payload = ujson.dumps(msg)
    deliver(payload)
    return text('', status=204)


@app.route('/api/status/<cluster>/<id:int>/<state:int>', methods=['POST'])
async def status(request, cluster: str, id: int, state: int):
    valid = validate_shard_in_cluster(cluster, id)
    if valid is not None: return valid
    
    State[cluster][str(id)] = {
        'state': state,
        'last_healthy': time.time()
    }
    payload = ujson.dumps({
        'op': 'state_update',
        'payload': {
            'cluster': cluster,
            'shard': id,
            'state': state
        }
    })
    deliver(payload)
    return text('', status=204)


@app.websocket('/api/feed')
async def feed(request, ws):
    Feeds.append(ws)
    ws.connection_lost = lambda exc: Feeds.remove(ws)
    hello = {
        'op': 'hello',
        'payload': State
    }
    await ws.send(ujson.dumps(hello))
    while True:
        # Keep the socket open indefinitely
        # TODO: is there a better solution to this
        await asyncio.sleep(15)


@app.route('/')
async def index(request):
    return await file('./index.html')


@app.route('/plain')
async def noscript(request):
    return text(str(State))


def deliver(payload):
    for client in Feeds:
        asyncio.ensure_future(client.send(payload))


@app.listener('after_server_start')
async def hook_health(app, loop):
    asyncio.ensure_future(_health_loop())
    logging.info('health loop running')


async def _health_loop():
    while True:
        print('running health checks')
        await _run_health_checks()
        print('health checks completed, sleeping')
        await asyncio.sleep(health_timeout)

async def _run_health_checks():
    now = time.time()
    for (cluster, shards) in State.items():
        for (id, shard) in shards.items():
            if id == 'len': continue
            
            if shard['state'] == ShardState.UNHEALTHY: # this shard is already sick
                continue # TODO: alerting for continued sickness?
            
            if (now - shard['last_healthy']) > 15:
                shard['state'] = ShardState.UNHEALTHY
                payload = ujson.dumps({
                    'op': 'state_update',
                    'payload': {
                        'cluster': cluster,
                        'shard': id,
                        'state': shard['state']
                    }
                })
                deliver(payload)


def configure_state():
    # GLANCE_BOTS=main:100;patron_1:6;patron_2:6
    bots_raw = os.getenv('GLANCE_BOTS')
    for bot in bots_raw.split(';'):
        (cluster, shards) = bot.split(':')
        shards = int(shards)
        State[cluster] = dict()
        State[cluster]['len'] = shards
    

if __name__ == '__main__':
    load_dotenv()
    configure_state()
    health_timeout = int(os.getenv('GLANCE_HEALTH_TIMEOUT') or '15')

    app.run(host=os.getenv('GLANCE_HOST') or '0.0.0.0',
            port=int(os.getenv('GLANCE_PORT') or '3000')
            )
