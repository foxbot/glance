from aiohttp import ClientSession
import asyncio
import random

base = 'http://localhost:3050'
limit = 260
shards = 280

async def bound_post(sem, url, sesion):
    async with sem:
        await session.post(url)

async def run_init(session):
    tasks = []
    sem = asyncio.Semaphore(limit)
    for i in range(shards):
        url = f'{base}/api/status/main/{i}/3'
        task = asyncio.ensure_future(bound_post(sem, url, session))
        tasks.append(task)
    await asyncio.gather(*tasks)

async def run_healthy(session):
    healthy = list(range(shards))
    sem = asyncio.Semaphore(limit)
    while True:
        if len(healthy) > 90:
            healthy.remove(random.choice(healthy))
        random.shuffle(healthy)
        tasks = []
        
        for s in healthy:
            url = f'{base}/api/health/main/{s}'
            task = asyncio.ensure_future(bound_post(sem, url, session))
            await asyncio.sleep(random.random() / 4)
            tasks.append(task)
        await asyncio.gather(*tasks)

async def run_demo():
    while True:
        shard = random.randint(0, shards-1)
        state = random.randint(0, 6)
        url = f'{base}/api/status/main/{shard}/{state}'
        del url

async def main(session):
    await run_init(session)
    await run_healthy(session)

loop = asyncio.get_event_loop()
with ClientSession() as session:
    loop.run_until_complete(main(session))
