# boilerplate-cli-ui-go-v2-react

Go CLI with embedded React 18 web UI. Single binary, no runtime dependencies.

Part of [SuperCLI](https://github.com/javimosch/supercli) - build CLI/UI plugins fast for 2026.

**v1**: [boilerplate-cli-ui-go](https://github.com/javimosch/boilerplate-cli-ui-go) | **Vue**: [boilerplate-cli-ui-go-v2-vue](https://github.com/javimosch/boilerplate-cli-ui-go-v2-vue) | **Node**: [boilerplate-cli-ui-node](https://github.com/javimosch/boilerplate-cli-ui-node)

## Architecture

```
boilerplate-cli-ui-go-v2-react/
├── main.go           # CLI entry point (start, stop, status, version)
├── server.go         # HTTP server with go:embed for UI files
├── daemon.go         # Daemon management (pid file, signals)
├── ui/               # Frontend (embedded at compile time)
│   ├── index.html    # Entry point (React 18 from CDN)
│   ├── css/
│   │   └── app.css
│   └── js/
│       ├── app.jsx           # React app with routing
│       ├── components/
│       │   ├── Sidebar.jsx
│       │   └── StatusCard.jsx
│       └── views/
│           ├── Dashboard.jsx
│           └── Settings.jsx
├── go.mod
├── build.sh
└── README.md
```

## Key Feature: `go:embed`

Frontend files are **separate** but **embedded into the binary** at compile time:

```go
//go:embed ui/*
var uiFiles embed.FS
```

**Benefits:**
- Single binary output (no runtime file dependencies)
- Separate HTML/CSS/JSX files (proper syntax highlighting)
- No build step for frontend (CDN-based React)
- Hot-reload during development (serve from disk)

## Hashbang Routing

Routes use hashbang URLs:
- `http://localhost:8080/#/dashboard` - Dashboard view
- `http://localhost:8080/#/settings` - Settings view
- `http://localhost:8080/` - Defaults to dashboard

## Build

```bash
chmod +x build.sh
./build.sh
```

Output: Single binary `boilerplate-cli-ui-go-v2-react`

## Usage

```bash
# Start server (foreground)
./boilerplate-cli-ui-go-v2-react start

# Start on custom port
./boilerplate-cli-ui-go-v2-react start -port 3000

# Start as daemon
./boilerplate-cli-ui-go-v2-react start -daemon

# Stop daemon
./boilerplate-cli-ui-go-v2-react stop

# Check status
./boilerplate-cli-ui-go-v2-react status
```

## API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /` | Web UI |
| `GET /api/status` | Server status (JSON) |
| `GET /api/health` | Health check (JSON) |

## Frontend Stack

- **React 18** (CDN) - UI library with Hooks
- **Tailwind CSS** (CDN) - Utility-first styling
- **Lucide Icons** (CDN) - Icon library
- **Babel Standalone** (CDN) - JSX transformation
- **Hashbang routing** - `#/dashboard`, `#/settings`

No npm, no build step. All code in `ui/index.html`.

## Development

### Option 1: Edit embedded files

1. Edit files in `ui/`
2. Run `go run .` (files are re-embedded each run)
3. Refresh browser

### Option 2: Serve from disk (faster)

For development, you can serve files directly from disk:

```go
// In server.go, temporarily replace:
// uiSub, _ := fs.Sub(uiFiles, "ui")
// fileServer := http.FileServer(http.FS(uiSub))

// With:
fileServer := http.FileServer(http.Dir("ui"))
```

This allows hot-reload without recompiling.

## Adding New Views

1. Create `ui/js/views/MyView.jsx`:

```jsx
function MyView() {
    const [data, setData] = React.useState(null);

    React.useEffect(() => {
        fetch('/api/my-data')
            .then(r => r.json())
            .then(setData);
    }, []);

    return (
        <div>
            <h2>My View</h2>
            {/* Your content */}
        </div>
    );
}
```

2. Add view in `ui/js/app.jsx`:

```jsx
// In the render section
{currentView === 'my-view' && <MyView />}
```

3. Add nav item in `ui/js/components/Sidebar.jsx`:

```jsx
const navItems = [
    // ... existing items
    { id: 'my-view', label: 'My View', icon: 'star' },
];
```

4. Include script in `ui/index.html`:

```html
<script type="text/babel" src="/js/views/MyView.jsx"></script>
```

## Adding New API Endpoints

1. Add handler in `server.go`:

```go
func handleMyEndpoint(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"hello": "world"})
}
```

2. Register in `startServer()`:

```go
mux.HandleFunc("/api/my-endpoint", handleMyEndpoint)
```

## Comparison with Vue Version

| Aspect | Vue 3 Version | React Version |
|--------|---------------|---------------|
| Framework | Vue 3 (CDN) | React 18 (CDN) |
| JSX | No (templates) | Yes (Babel) |
| Routing | Manual state | Manual state |
| State | provide/inject | Context API |
| File extension | .js | .jsx |
| Learning curve | Lower | Higher |

Both versions:
- Single binary output
- Embedded frontend via `go:embed`
- CDN-based (no npm)
- Tailwind CSS + Lucide Icons
