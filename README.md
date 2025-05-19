# Atuin Web UI

A web-based UI for [Atuin](https://github.com/atuinsh/atuin) shell history, with the backend serving the frontend in an embedded mode.

## Features

- View and search your Atuin shell history
- Filter history by various criteria
- Delete history entries
- Cross-platform support (Windows, macOS, Linux)
- Single binary distribution with embedded frontend

## Building from Source

### Prerequisites

- Go 1.16 or later
- Node.js 14 or later
- npm

> **Note:** This application uses a pure Go implementation of SQLite (modernc.org/sqlite) that doesn't require CGO or a C compiler.

### Building

#### Windows

```powershell
# Run the build script
.\build.ps1
```

#### macOS/Linux

```bash
# Make the build script executable
chmod +x build.sh

# Run the build script
./build.sh
```

The build script will:
1. Build the Vue.js frontend
2. Embed the frontend files into the Go backend
3. Compile the Go backend with the embedded frontend
4. Create binaries for Windows, macOS, and Linux in the `dist` directory

> **Note about cross-compilation:** This application uses a pure Go implementation of SQLite that doesn't require CGO, allowing for easy cross-compilation to all platforms from any development environment.

## Development

### Frontend

```bash
cd frontend
npm install
npm run serve
```

The frontend development server will run at http://localhost:8081

### Backend

```bash
cd backend
go run .
```

The backend API server will run at http://localhost:8080

## Usage

1. Download the appropriate binary for your platform from the releases page or build from source
2. Run the binary
3. Access the UI in your web browser at http://localhost:8080

The application will attempt to automatically locate your Atuin database. If it cannot find it, you'll be prompted to specify the path.

## Configuration

The application stores its configuration in:
- Windows: `%APPDATA%\atuin-ui\config.json`
- macOS: `~/Library/Application Support/atuin-ui/config.json`
- Linux: `~/.config/atuin-ui/config.json`

## License

[MIT License](LICENSE)
