from requests_threads import AsyncSession

session = AsyncSession(n=40)


BASE = 'http://localhost:3050'
SHARDS = 112
BOT = 'main'

async def _main():
    # populate online initially
    for x in range(SHARDS):
        session.post(f'{BASE}/api/status/{BOT}/{x}/2')

session.run(_main)


print('ok')
