# glance
The at-a-glance Discord bot status dashboard.

## Features

Glance is intended to be used by large, multi-process Discord bots, where keeping
tabs on bot status becomes more complicated than just checking the user's status
on Discord.

The simple user-interface provided by Glance makes it a compelling status page to
offer as part of your bot's public status page, allowing your end users to quickly
check the bot's status in their guild.

Shards are shown in a grid, with the status of the shard being indicated by universally
understood traffic-light colors<sup>[1](https://pdfs.semanticscholar.org/738b/fe8606a556e0a1fe85686c5c20616a1013dd.pdf)</sup>.

A search bar is provided at the top of the page where a user may enter their Guild ID.
Glance will then locate that guild in each bot cluster using Discord's standard guild
selection algorithm, and highlight the shard in the list. This allows users to easily
locate the status of the bot in a multi-sharded, multi-cluster configuration.

![Glance Screenshot](https://i.foxbot.me/9v5ijLF.png)

### Non-features

Glance provides a very basic level of data-collection using push-based metrics. This
allows the status page to determine the live status of the shard, and mark a shard as
unhealthy when it has not been regularly publishing its status.

The limited nature of Glance's data collection means that Glance is not a replacement
for a more complete data-visualization package, such as Grafana. While Glance can be
very helpful for your bot's support team, or even for your developers/ops team to reference,
it does not replace the infinitely more powerful visualization and alerting that a
time-series database paired with Grafana is capable of.

## Installation

Glance requires the following:

- Python 3.7
- Optionally, an ASGI-compatible Web Server

`$ pip install -r requirements.txt`

### Configuration

Copy the `.env.sample` file to `.env`, or, export the environment variables yourself.

See [`.env.sample`](./.env.sample) for documentation on how to set configuration flags. Additionally,
you may [configure the integrated server](https://sanic.readthedocs.io/en/latest/sanic/config.html#builtin-configuration-values)
with environment variables.

### Deployment

Glance is composed of a Python web-server and a single static-page site. The static-page, located
at `./index.html` and `assets/` can be served directly, with minimal changes (see below). The API
must be served through a WSGI-compatible web-server, or using the out-of-box integrated server.

#### Interated Server

`$ python3 app.py`

Glance and its API will be accessible at the address specified in the configuration.

#### WSGI/ASGI Server

See [Sanic's documentation](https://sanic.readthedocs.io/en/latest/sanic/deploying.html#running-via-asgi)
on this subject.

Yes, WSGI through Gunicorn is supported.

**Note:** you may need to configure the frontend differently when using an external server,
as it will try to find the socket from a fixed location. The `noscript` link also assumes
use of the integrated server, if that's an edge case you care to fix.

## Usage

Bots wishing to push data to Glance must post data to the Glance API both when the shard
state is updated, and at (around) the specified health interval.

By default, the health interval is set to 30 seconds (i.e., the bot must contact Glance every
30 seconds, lest it be considered unhealthy). It may make more sense in your environment
to configure this value as 45 seconds, and post a health check every heartbeat.

You may also want to account for other data before sending a health check, such as whether
or not the shard is receiving a normal amount of messages from the gateway. This might help
you to identify if a shard is in a zombie state (i.e., your library has deadlocked, or has lost
track of its connection state) and display a status update accordingly.

Note: If a shard is marked as unhealthy, and Glance receives a status update for it, the shard 
will then be marked as ONLINE; not whichever state it was in prior to going unhealthy. This may
be something to consider before sending health checks.

### Frontend Configuration

You are free<sup>*</sup> to modify the header of `index.html` to include your branding (the intended location
is where `foxbot.me` is located in `.header-author`, though you can place it anywhere). 

Class names and element IDs are how the site's pure-vanilla scripting works, so be aware that changing
them will likely break the site.

All colors and fonts used on the frontend may be adjusted via `assets/main.css`. I don't load any fonts
from the internet by default, so you can work that out on your own.

**If you are running Glance outside of the integrated server**, you will need to adjust the `getSocketUrl`
method under `assets/index.js`. The frontend will not be able to connect to the API if you don't
adjust this.

<small>*"free" as in "encouraged". This statement carries no implications on licensing; please see the
`LICENSE` for licensing guarantees.</small>

### API Documentation

If Glance is configured with `GLANCE_PASSPHRASE`, requests to the following endpoints must have
their `Authorization` header set accordingly.

##### Shard State Enumeration
| Status Name | Value |
| --- | --- |
| Unknown | 0 |
| Unhealthy | 1 |
| Offline | 2 |
| Online | 3 |
| Starting | 4 |
| Stopping | 5 |
| Resuming | 6 |

Note: shard states generally have no real meaning, and are generally only used to
set the color of the shard on the Glance grid.

The unknown, unhealthy, and online states are used internally by Glance. Shards are
initialized with the unknown state, will toggle to unhealthy when the health check
interval has been exceeded, and will toggle back to online if a health check is
received for an unhealthy shard.

#### `POST /api/status/<cluster>/<id:int>/<state:int>`
This endpoint sets the state for a given shard.

##### Parameters
`cluster`: The name of the cluster this shard is on; this should match up with the cluster names in the configuration.

`id`: The zero-indexed ID of this shard.

`state`: The numeric Shard State of this shard; see above.

#### `POST /api/health/<cluster>/<id:int>`
This endpoint marks a shard as healthy.

##### Parameters
`cluster`: The name of the cluster this shard is on; this should match up with the cluster names in the configuration.

`id`: The zero-indexed ID of this shard.

## License
Glance is licensed under ISC. See the [LICENSE](./LICENSE.md), located at `LICENSE.md`

## Contributing
1. Please open an issue or contact me before opening a merge request.
2. Adhere to the standards set in the `.editorconfig`.
3. All changes should make the code better in some way. See http://suckless.org/philosophy/
