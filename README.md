# gostatus 

[WIP] Discord presence badges as SVGs to show status, Spotify or game activity using [shields.io](https://shields.io/). *Inspired by [statusbadges](https://github.com/advaith1/statusbadges), rewritten in Go.*

## Running

### Requirements:
- [Go 1.21+](https://go.dev/)
- **Server Members** and **Presence** intents need to be enabled at [Dev Portal](https://discord.com/developers/home)
- Bot must be in a shared server with the users you want to track

```bash
git clone https://github.com/vmphase/gostatus
cd gostatus
go mod tidy

# Linux
DISCORD_TOKEN=your_token_here go run main.go

# Windows
$env:DISCORD_TOKEN="your_token_here"; go run main.go
``` 

The server starts on port `8080` by default. Set `PORT` environmental variable to override.


## Endpoints

### `GET /badge/status/{discord_user_id}`
Returns a badge with the user's current Discord status.

| Query        | Default      | Description |
| -----        | -------      | ----------- |
| `label`      | `currently`  | Left side text |
| `color`      | status-based | Right side background color |
| `labelColor` | `#555`       | Left side background color |
| `simple`     | —            | Set to `true` to collapse `idle`/`dnd` => `online` |

---

### `GET /badge/spotify/{discord_user_id}`
Returns a badge with the track the user is currently listening to on Spotify.

| Query        | Default        | Description |
| -----        | -------        | ----------- |
| `label`      | `listening to` | Left side text |
| `color`      | `#1db954`      | Right side background color |
| `labelColor` | `#555`         | Left side background color |
| `fallback`   | `nothing`      | Text shown when not listening |

---

### `GET /badge/playing/{discord_user_id}`
Returns a badge with the game the user is currently playing.

| Query        | Default      | Description |
| ------       | -------      | ----------- |
| `label`      | `playing`    | Left side text |
| `color`      | `#5865f2`    | Right side background color |
| `labelColor` | `#555`       | Left side background color |
| `fallback`   | `nothing`    | Text shown when not playing |

---

### `GET /presence/{discord_user_id}`
Returns raw presence data as JSON. CORS-enabled.
