# gostatus 

Discord presence badges as SVGs to show discord status, spotify-, vscode-, zed- or game-activity using [shields.io](https://shields.io/).   
> [!NOTE]
> Originally inspired by [statusbadges](https://github.com/advaith1/statusbadges), rewritten in Go with improvments.

## Running

### Requirements:
- [Go 1.21+](https://go.dev/)
- **Server Members** and **Presence** intents need to be enabled at [Dev Portal](https://discord.com/developers/home)
- Bot must be in a shared server with the users you want to track

```bash
git clone https://github.com/vmphase/gostatus
cd gostatus

go mod tidy
go run main.go
``` 

Configuration is done via `config.[example].toml`:
- rename `config.toml.example` to `config.toml`
- set `token` to your bot token
- optionally change the `port` field

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
| `hideLogo`   | `false`        | Set to `true` to hide the Spotify logo |

---

### `GET /badge/crunchyroll/{discord_user_id}`

Returns a badge with the information the user is currently watching on Crunchyroll.
| Query        | Default        | Description |
| -----        | -------        | ----------- |
| `label`      | `watching`     | Left side text |
| `color`      | `#5865f2`      | Right side background color |
| `labelColor` | `#555`         | Left side background color |
| `fallback`   | `nothing`      | Text shown when not listening |
| `hideLogo`   | `false`        | Set to `true` to hide the Crunchyroll logo |

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

### `GET /badge/vscode/{discord_user_id}`

Returns a badge with the file and workspace the user is currently editing in [Visual Studio Code](https://code.visualstudio.com/).
| Query        | Default      | Description |
| ------       | -------      | ----------- |
| `label`      | `vscode`     | Left side text |
| `color`      | `#23a7f2`    | Right side background color |
| `labelColor` | `#555`       | Left side background color |
| `fallback`   | `nothing`    | Text shown when not coding |
| `hideLogo`   | `false`      | Set to `true` to hide the VSCode logo |

---

### `GET /badge/zed/{discord_user_id}`

Returns a badge with the file and workspace the user is currently editing in [Zed](https://zed.dev/).
| Query        | Default      | Description |
| ------       | -------      | ----------- |
| `label`      | `zed`        | Left side text |
| `color`      | `#a277ff`    | Right side background color |
| `labelColor` | `#2b213a`    | Left side background color |
| `fallback`   | `nothing`    | Text shown when not coding |
| `hideLogo`   | `false`      | Set to `true` to hide the Zed logo |

---

### `GET /presence/{discord_user_id}`
Returns raw presence data as JSON. CORS-enabled.
 