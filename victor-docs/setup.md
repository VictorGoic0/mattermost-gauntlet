# Mattermost Local Development Setup Guide

This guide will walk you through setting up Mattermost for local development on your machine.

## Prerequisites

### Required Software

1. **Go** (1.21 or later)
   - Download: https://go.dev/dl/
   - Verify: `go version`

2. **Node.js** (20.x LTS) and npm
   - Download: https://nodejs.org/
   - Verify: `node --version` and `npm --version`

3. **Docker Desktop** (for database)
   - Download: https://www.docker.com/products/docker-desktop/
   - Verify: `docker --version` and `docker compose version`

4. **Git**
   - Download: https://git-scm.com/downloads
   - Verify: `git --version`

### Recommended Tools

- **Make** (usually pre-installed on macOS/Linux, Windows users need to install)
  - Windows: Install via Chocolatey: `choco install make`
- **Code Editor**: VS Code or Cursor
- **PostgreSQL Client** (optional, for database inspection): Postico, TablePlus, or pgAdmin

---

## Step 1: Clone the Repository
```bash
# Clone the main Mattermost repository
git clone https://github.com/mattermost/mattermost.git
cd mattermost

# Create your feature branch
git checkout -b your-feature-name
```

**Repository Structure**:
```
mattermost/
├── server/          # Go backend code
├── webapp/          # React frontend code
├── e2e-tests/       # End-to-end tests
└── docker-compose.yaml
```

---

## Step 2: Start the Database

Mattermost requires PostgreSQL (or MySQL, but we'll use Postgres).

### Option A: Using Docker Compose (Recommended)
```bash
# Start PostgreSQL in Docker
docker compose up -d postgres

# Verify it's running
docker compose ps
```

This starts PostgreSQL on `localhost:5432` with:
- Database: `mattermost_test`
- User: `mmuser`
- Password: `mostest`

### Option B: Local PostgreSQL Installation

If you prefer a local Postgres install:
```bash
# macOS (via Homebrew)
brew install postgresql@14
brew services start postgresql@14

# Create database and user
psql postgres
CREATE DATABASE mattermost_test;
CREATE USER mmuser WITH PASSWORD 'mostest';
GRANT ALL PRIVILEGES ON DATABASE mattermost_test TO mmuser;
\q
```

---

## Step 3: Configure the Server

### Create Config File
```bash
cd server

# Copy the default config
cp config/config.json config/config-local.json
```

### Edit `config/config-local.json`

Key settings to verify/modify:
```json
{
  "ServiceSettings": {
    "SiteURL": "http://localhost:8065",
    "ListenAddress": ":8065",
    "EnableDeveloper": true
  },
  "SqlSettings": {
    "DriverName": "postgres",
    "DataSource": "postgres://mmuser:mostest@localhost:5432/mattermost_test?sslmode=disable&connect_timeout=10"
  },
  "FileSettings": {
    "Directory": "./data/"
  },
  "PluginSettings": {
    "Enable": true,
    "EnableUploads": true,
    "Directory": "./plugins",
    "ClientDirectory": "./client/plugins"
  }
}
```

**Important Settings**:
- `EnableDeveloper`: Set to `true` for development mode
- `DataSource`: Update if you're not using Docker Compose defaults
- `PluginSettings.Enable`: Must be `true` if building plugins

---

## Step 4: Build and Run the Server
```bash
# From the server/ directory

# Install Go dependencies
go mod download

# Build the server
make build-server

# Run database migrations
make run-server migrate

# Start the server
make run-server

# OR use this single command that does it all:
make run
```

The server will start on `http://localhost:8065`

**Expected Output**:
```
{"level":"info","msg":"Server is listening on :8065"}
{"level":"info","msg":"Starting Server..."}
```

### Troubleshooting Server Issues

**Port already in use**:
```bash
# Find what's using port 8065
lsof -i :8065

# Kill the process
kill -9 <PID>
```

**Database connection errors**:
```bash
# Check if Postgres is running
docker compose ps

# Check logs
docker compose logs postgres

# Restart Postgres
docker compose restart postgres
```

---

## Step 5: Build and Run the Web App

Open a **new terminal window** (keep the server running in the first).
```bash
# From the project root
cd webapp

# Install npm dependencies (this takes a few minutes)
npm install

# Start the development server
npm run dev
```

The web app will start on `http://localhost:8065` and proxy API requests to the server.

**Expected Output**:
```
VITE v4.x.x ready in XXX ms
➜ Local:   http://localhost:8065/
```

---

## Step 6: Create Your First Account

1. Open browser to `http://localhost:8065`
2. Click **"Create an account"**
3. Fill in:
   - Email: `admin@example.com`
   - Username: `admin`
   - Password: `password123`
4. Click **"Create Account"**
5. You'll be prompted to create a team - create one (e.g., "Dev Team")

**First user is automatically a System Admin** with full permissions.

---

## Step 7: Verify the Setup

### Check Server Health
```bash
curl http://localhost:8065/api/v4/system/ping
# Expected: {"status":"OK"}
```

### Check Database Connection
```bash
# Connect to Postgres
docker compose exec postgres psql -U mmuser -d mattermost_test

# List tables
\dt

# Should see tables like: users, channels, posts, teams, etc.

# Exit
\q
```

### Check Plugin Directory
```bash
cd server
ls -la plugins/
# Should exist (may be empty initially)
```

---

## Step 8: Development Workflow

### Making Changes

**Backend (Go) Changes**:
```bash
# Server auto-reloads with 'make run-server'
# If not, restart manually:
cd server
make restart-server
```

**Frontend (React) Changes**:
```bash
# Webpack dev server auto-reloads
# If you see issues, restart:
cd webapp
npm run dev
```

### Running Tests

**Server Tests**:
```bash
cd server
make test-server
```

**Web App Tests**:
```bash
cd webapp
npm run test
```

---

## Development Tools

### Useful Make Commands

From `server/` directory:
```bash
make run               # Build and run server
make stop-server       # Stop the server
make clean             # Clean build artifacts
make build-server      # Build server binary
make test-server       # Run server tests
make govet             # Run Go vet
make check-style       # Check code style
```

From `webapp/` directory:
```bash
npm run dev            # Start dev server
npm run build          # Production build
npm run test           # Run tests
npm run lint           # Lint code
npm run check-types    # TypeScript type checking
```

### Recommended VS Code / Cursor Extensions

- **Go** (`golang.go`)
- **ESLint** (`dbaeumer.vscode-eslint`)
- **Prettier** (`esbenp.prettier-vscode`)
- **GitLens** (`eamodio.gitlens`)

---

## Plugin Development Setup

If you're building a plugin:

### 1. Enable Plugin Support

Already done if you followed Step 3 (config.json has `PluginSettings.Enable: true`)

### 2. Plugin Development Mode
```bash
cd server

# Set environment variable
export MM_PLUGINSETTINGS_ENABLEUPLOADS=true
```

### 3. Create Plugin Structure
```bash
# From server/ directory
mkdir -p plugins/your-plugin-name
cd plugins/your-plugin-name

# Create plugin manifest
cat > plugin.json << 'EOF'
{
  "id": "com.example.yourplugin",
  "name": "Your Plugin Name",
  "version": "0.1.0",
  "min_server_version": "9.0.0",
  "server": {
    "executable": "server"
  },
  "settings_schema": {}
}
EOF
```

### 4. Build and Install Plugin
```bash
# Build plugin
make build

# Upload to Mattermost (via UI or API)
# System Console > Plugins > Upload Plugin
```

**Plugin Documentation**: https://developers.mattermost.com/integrate/plugins/

---

## Common Issues and Solutions

### Issue: `make: command not found`

**Solution**: Install Make
```bash
# macOS
xcode-select --install

# Windows (via Chocolatey)
choco install make

# Linux (Ubuntu/Debian)
sudo apt-get install build-essential
```

### Issue: Go build errors about missing modules

**Solution**: 
```bash
cd server
go mod tidy
go mod download
```

### Issue: npm install fails

**Solution**:
```bash
cd webapp
rm -rf node_modules package-lock.json
npm cache clean --force
npm install
```

### Issue: "Port 8065 already in use"

**Solution**:
```bash
# Find the process
lsof -i :8065

# Kill it
kill -9 <PID>

# Or change the port in config.json
```

### Issue: Database migration errors

**Solution**:
```bash
# Reset database
docker compose down -v
docker compose up -d postgres

# Re-run migrations
cd server
make run-server migrate
```

---

## Environment Variables Reference
```bash
# Database
export MM_SQLSETTINGS_DATASOURCE="postgres://mmuser:mostest@localhost/mattermost_test?sslmode=disable"

# Server
export MM_SERVICESETTINGS_SITEURL="http://localhost:8065"
export MM_SERVICESETTINGS_ENABLEDEVELOPER=true

# Plugins
export MM_PLUGINSETTINGS_ENABLE=true
export MM_PLUGINSETTINGS_ENABLEUPLOADS=true

# Logging
export MM_LOGSETTINGS_CONSOLELEVEL=DEBUG
```

---

## Next Steps

1. **Read the Developer Documentation**: https://developers.mattermost.com/
2. **Explore the Codebase**: Start with `server/app/` and `webapp/channels/src/`
3. **Join the Community**: https://community.mattermost.com/
4. **Review Contributing Guidelines**: https://developers.mattermost.com/contribute/

---

## Useful Links

- **Main Docs**: https://docs.mattermost.com/
- **Developer Docs**: https://developers.mattermost.com/
- **API Docs**: https://api.mattermost.com/
- **Plugin Docs**: https://developers.mattermost.com/integrate/plugins/
- **Community Forum**: https://community.mattermost.com/
- **GitHub Issues**: https://github.com/mattermost/mattermost/issues

---

## Quick Reference Commands
```bash
# Start everything
docker compose up -d postgres
cd server && make run &
cd webapp && npm run dev

# Stop everything
pkill -f mattermost
docker compose down

# Reset everything
docker compose down -v
rm -rf server/data
dropdb mattermost_test && createdb mattermost_test
```

---

**You're all set!** 🚀

If you run into issues, check the [troubleshooting section](#common-issues-and-solutions) or ask in the Mattermost community forums.