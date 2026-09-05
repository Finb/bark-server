# Bark

<img src="https://wx3.sinaimg.cn/mw690/0060lm7Tly1g0nfnjjxbbj30sg0sg757.jpg" width=200px height=200px />

[Bark](https://github.com/Finb/Bark) is an iOS App which allows you to push customed notifications to your iPhone.

## Installation

### For Docker User

![Docker Automated build](https://img.shields.io/docker/automated/finab/bark-server.svg) ![Image Size](https://img.shields.io/docker/image-size/finab/bark-server?sort=date) ![License](https://img.shields.io/github/license/finb/bark-server)

The docker image is already available, you can use the following command to run the bark server:

``` sh
docker run -dt --name bark -p 8080:8080 -v `pwd`/bark-data:/data finab/bark-server
```

You can also use the GitHub Container Registry image:

``` sh
docker run -dt --name bark -p 8080:8080 -v `pwd`/bark-data:/data ghcr.io/finb/bark-server
```

If you use the docker-compose tool, you can copy docker-copose.yaml under this project to any directory and run it:

``` sh
mkdir bark-server && cd bark-server
curl -sL https://github.com/Finb/bark-server/raw/master/deploy/docker-compose.yaml > docker-compose.yaml
docker compose up -d
```

### For General User 

- 1、Download precompiled binaries from the [releases](https://github.com/Finb/bark-server/releases) page
- 2、Add executable permissions to the bark-server binary: `chmod +x bark-server`
- 3、Start bark-server: `./bark-server --addr 0.0.0.0:8080 --data ./bark-data`
- 4、Test the server: `curl localhost:8080/ping`

**Note: Bark-server uses the `/data` directory to store data by default. Make sure that bark-server has permission to write to the `/data` directory, otherwise use the `-d` option to change the directory.**

### For Developer

Developers can compile this project by themselves, and the dependencies required for compilation:

- Golang 1.18+
- Go Mod Enabled(env `GO111MODULE=on`)
- Go Mod Proxy Enabled(env `GOPROXY=https://goproxy.cn`)
- [go-task](https://taskfile.dev/installation/) Installed

Run the following command to compile this project:

```sh
# Cross compile all platforms
task

# Compile the specified platform (please refer to Taskfile.yaml)
task linux_amd64
task linux_amd64_v3
```

**Note: The linux amd64 v3 architecture was added in go 1.18, see [https://github.com/golang/go/wiki/MinimumRequirements#amd64](https://github.com/golang/go/wiki/MinimumRequirements#amd64)**

### Use MySQL instead of Bbolt

Just run the server with `-dsn=user:pass@tcp(mysql_host)/bark`, it will use MySQL instead of file database Bbolt

## Others

* [API_V2.md](docs/API_V2.md).
* [MCP.md](docs/MCP.md).



## 🌐 Web Resources & Verified Articles Directory
- [COLOR BRAIN TEST GAMES](https://pixelarcade-speed.web.app/color-brain-test-games.html)
- [TILE SORT MATCH 3](https://koreagame-hub24.netlify.app/tile-sort-match-3.html)
- [ONLY UP](https://quantum-puzzle-hub.pages.dev/only-up.html)
- [BRAIN TEST ONE LINE DRAW PUZZLE](https://congdonggame-vietnam.web.app/brain-test-one-line-draw-puzzle.html)
- [REFLECT BEAM LASER LOGIC](https://francejeux-online.web.app/reflect-beam-laser-logic.html)
- [DIG OUT OF PRISON](https://bharat-game-zone.web.app/dig-out-of-prison.html)
- [GIRL RESCUE DRAGON OUT](https://jogosonline-brasil.vercel.app/girl-rescue-dragon-out.html)
- [STEAL BRAINROT MONSTERS](https://kuaile-youxi-hub.web.app/steal-brainrot-monsters.html)
- [FRUIT CONNECT 3](https://dautruong-game24h.web.app/fruit-connect-3.html)
- [DUCK LUCK](https://luchshie-igry-rus.pages.dev/duck-luck.html)
- [THE WALL](https://action-strike-zone.pages.dev/the-wall.html)
- [BARBEE BLACK FRIDAY FASHION](https://zona-igr-besplatno.web.app/barbee-black-friday-fashion.html)
- [BLOCK UP](https://webarcade-hub.github.io/block-up.html)
- [EGG DASH](https://pixelarcade-speed.web.app/egg-dash.html)
- [M5 CITY DRIVER](https://jogosonline-brasil.vercel.app/m5-city-driver.html)
- [SURVIVE LAVA FOR BRAINROTS](https://pixelarcadezgame.web.app/survive-lava-for-brainrots.html)
- [GUESS WORD](https://luchshie-igry-rus.pages.dev/guess-word.html)
- [STACK N SORT](https://mir-igr-onlayn.pages.dev/stack-n-sort.html)
- [HELP ME TRICKY BRAIN PUZZLES](https://juegosweb-desbloqueados.vercel.app/help-me-tricky-brain-puzzles.html)
- [SUPERMARKET MANAGER SIMULATOR](https://jogosweb-brasil.github.io/supermarket-manager-simulator.html)
- [UPHILL RUSH 13](https://koreagame-arcade.netlify.app/uphill-rush-13.html)
- [MAHJONG RIDDLES EGYPT](https://unblocked-galaxy-hub.pages.dev/mahjong-riddles-egypt.html)
- [ANIME WOW](https://instantsounds-daw.pages.dev/sound/anime-wow.html)
- [TRIAL XTREME](https://webarcade-gamehub.github.io/trial-xtreme.html)
- [DREAM ROOM MAKEOVER](https://unblocked-galaxy-hub.pages.dev/dream-room-makeover.html)
- [ZINDEX](https://choigame24h-vietnam.netlify.app/zindex.html)
- [MAHJONG TILE CLUB](https://webarcade-gamehub.github.io/mahjong-tile-club.html)
- [THE KULKA](https://zona-juegos-flash.web.app/the-kulka.html)
- [DRAGON](https://blox-trade-fairness.pages.dev/values/dragon)
- [CITY TOWER BUILDER](https://youxi-china24.netlify.app/city-tower-builder.html)
- [CUTE COLORING GAMES](https://hindigames-portal.netlify.app/cute-coloring-games.html)
- [SUPERMARKET SIMULATOR DREAM STORE](https://congdonggame-vietnam.web.app/supermarket-simulator-dream-store.html)
- [DUO FAMILY SANTA](https://planetejeux-france.pages.dev/duo-family-santa.html)
- [THE BODYGUARD](https://seoul-game-hub.pages.dev/the-bodyguard.html)
- [RACING MASTER 3D](https://mir-igr-onlayn.pages.dev/racing-master-3d.html)
- [BLOCK PUZZLE 3D](https://kuaile-youxi-hub.web.app/block-puzzle-3d.html)
- [ACCURATE 2D](https://onlinerus-portal.netlify.app/accurate-2d.html)
- [AIRPORT SECURITY](https://tokyo-arcade-web.pages.dev/airport-security.html)
- [SUGAR POP LAND](https://congdonggame-vietnam.web.app/sugar-pop-land.html)
- [MASTER BLENDER](https://arcadevault-gamehub.github.io/master-blender.html)
