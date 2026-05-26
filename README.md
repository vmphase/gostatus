# gostatus

Discord presence badges as SVGs to show discord status, music-, editor- or game-activity using [shields.io](https://shields.io/).

> [!NOTE]
> Originally inspired by [statusbadges](https://github.com/advaith1/statusbadges), rewritten in Go with improvements.

## Running

### Requirements

- [Go 1.21+](https://go.dev/)
- **Server Members** and **Presence** intents enabled at the [Dev Portal](https://discord.com/developers/home)
- Bot must be in a shared server with the users you want to track

### Local

```bash
git clone https://github.com/vmphase/gostatus
cd gostatus
go mod tidy
```

Copy the example config and set your bot token:

```bash
cp config.toml.example config.toml
```

```toml
token = "your_bot_token_here"
```

Then run:

```bash
go run main.go
```

The server listens on `:8080` by default. To change the port, update the `port` constant in `main.go`.

### Docker

#### Requirements

- [Docker](https://docs.docker.com/get-docker/) with Compose

#### Setup

Copy and configure:

```bash
cp config.toml.example config.toml
```

```toml
token = "your_bot_token_here"
```

Build and start:

```bash
docker compose -f .devcontainer/compose.yaml up --build -d
```

#### Changing the port

The port is hardcoded to `8080`. To change it, update it in three places:

1. `main.go` — the `port` constant
2. `.devcontainer/compose.yaml` — the `ports` mapping and `healthcheck` URL
3. `.devcontainer/Dockerfile` — the `EXPOSE` directive

## Endpoints

### `GET /badge/status/{discord_user_id}`

Current Discord presence status.

| Query        | Default      | Description                                         |
| ------------ | ------------ | --------------------------------------------------- |
| `label`      | `currently`  | Left side text                                      |
| `color`      | status-based | Right side background color                         |
| `labelColor` | `#555`       | Left side background color                          |
| `style`      | `flat`       | Badge style: `flat`, `flat-square`, `for-the-badge` |
| `simple`     | —            | Set to `true` to collapse `idle`/`dnd` => `online`  |

---

### `GET /badge/music/{discord_user_id}`

Track the user is currently listening to. Auto-detects supported music services (currently Spotify).

| Query        | Default        | Description                                         |
| ------------ | -------------- | --------------------------------------------------- |
| `label`      | `listening to` | Left side text                                      |
| `color`      | service-based  | Right side background color                         |
| `labelColor` | `#555`         | Left side background color                          |
| `style`      | `flat`         | Badge style: `flat`, `flat-square`, `for-the-badge` |
| `fallback`   | `nothing`      | Text shown when not listening                       |
| `hideLogo`   | `false`        | Set to `true` to hide the service logo              |

---

### `GET /badge/code/{discord_user_id}`

File and workspace the user is currently editing. Auto-detects supported editors (VSCode, Zed, Visual Studio).

| Query        | Default      | Description                                                                      |
| ------------ | ------------ | -------------------------------------------------------------------------------- |
| `label`      | editor-based | Left side text                                                                   |
| `color`      | editor-based | Right side background color                                                      |
| `labelColor` | `#1e1e2e`    | Left side background color                                                       |
| `style`      | `flat`       | Badge style: `flat`, `flat-square`, `for-the-badge`                              |
| `fallback`   | `nothing`    | Text shown when not coding                                                       |
| `hideLogo`   | `false`      | Set to `true` to hide the editor logo                                            |
| `prefer`     | first active | Preferred editor slug (`vscode`, `zed`, `visualstudio`) when multiple are active |

---

### `GET /badge/playing/{discord_user_id}`

Game the user is currently playing (editor activities are excluded).

| Query        | Default   | Description                                         |
| ------------ | --------- | --------------------------------------------------- |
| `label`      | `playing` | Left side text                                      |
| `color`      | `#5865f2` | Right side background color                         |
| `labelColor` | `#555`    | Left side background color                          |
| `style`      | `flat`    | Badge style: `flat`, `flat-square`, `for-the-badge` |
| `fallback`   | `nothing` | Text shown when not playing                         |

---

### `GET /badge/crunchyroll/{discord_user_id}`

Episode and series the user is currently watching on [Crunchyroll](https://www.crunchyroll.com/).

| Query        | Default    | Description                                         |
| ------------ | ---------- | --------------------------------------------------- |
| `label`      | `watching` | Left side text                                      |
| `color`      | `#f47521`  | Right side background color                         |
| `labelColor` | `#555`     | Left side background color                          |
| `style`      | `flat`     | Badge style: `flat`, `flat-square`, `for-the-badge` |
| `fallback`   | `nothing`  | Text shown when not watching                        |
| `hideLogo`   | `false`    | Set to `true` to hide the Crunchyroll logo          |

---

### `GET /presence/{discord_user_id}`

Raw presence data as JSON. CORS-enabled.
