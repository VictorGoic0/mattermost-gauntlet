# Mattermost Local Development Setup Guide

This guide will walk you through setting up Mattermost for local development on your machine.

## Prerequisites

### Required Software

1. **Go** (1.21 or later)
   - Download: https://go.dev/dl/
   - Verify: `go version`

2. **Node.js** (via NVM - recommended)
   - Install NVM: https://github.com/nvm-sh/nvm#installing-and-updating
   - Then install Node.js from within the `webapp` directory: `cd webapp && nvm install`
   - This ensures you get the correct Node.js version for Mattermost
   - Verify: `node --version` and `npm --version`
   - **Note for zsh users**: If you get `zsh: command not found: nvm`, add to `~/.zshrc`:
     ```bash
     export NVM_DIR="$HOME/.nvm"
     [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
     [ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"
     ```

3. **Docker Desktop** (for database and dependencies)
   - Download: https://www.docker.com/products/docker-desktop/
   - Verify: `docker --version` and `docker compose version`
   - **Windows users**: Use WSL2. Install via PowerShell (as admin): `wsl --install`
   - Make sure Docker has virtual file share access to your repository directory

4. **Git**
   - Download: https://git-scm.com/downloads
   - Verify: `git --version`

5. **Make**
   - macOS/Linux: Usually pre-installed
   - Ubuntu/Debian: `sudo apt install build-essential`
   - Windows: Install via Chocolatey: `choco install make` (or use WSL)

### Recommended Tools

- **Code Editor**: VS Code or Cursor
- **PostgreSQL Client** (optional, for database inspection): Postico, TablePlus, or pgAdmin

### System Configuration

**Increase file descriptors limit** (required):
```bash
# Add to your shell initialization file (~/.bashrc, ~/.zshrc, etc.)
ulimit -n 8096
```

**Install libpng** (required):
```bash
# macOS (via Homebrew)
brew install libpng

# Ubuntu/Debian
sudo apt-get install libpng-dev

# ARM-based Mac users: Install Rosetta (required for libpng)
softwareupdate --install-rosetta
```

---

## Step 1: Fork and Clone the Repository

1. **Fork the repository** on GitHub: https://github.com/mattermost/mattermost

2. **Clone your fork**:
```bash
git clone https://github.com/YOUR_GITHUB_USERNAME/mattermost.git
cd mattermost

# Create your feature branch
git checkout -b your-feature-name
```

**Note**: If you're creating a derivative version, you must comply with AGPLv2 license requirements and replace Mattermost branding per the trademark policy.

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
cd server

# Start the server (this will also start Docker dependencies)
make run-server
```

The server will start on `http://localhost:8065` and automatically:
- Start Docker containers (postgres, redis, etc.)
- Build the server
- Run database migrations
- Start the server process

**Expected Output**:
```
{"level":"info","msg":"Server is listening on :8065"}
{"level":"info","msg":"Starting Server..."}
```

**Test your environment**:
```bash
curl http://localhost:8065/api/v4/system/ping
```

Expected response:
```json
{"AndroidLatestVersion":"","AndroidMinVersion":"","DesktopLatestVersion":"","DesktopMinVersion":"","IosLatestVersion":"","IosMinVersion":"","status":"OK"}
```

**Set up admin user** (optional, but recommended):
```bash
# From the project root
bin/mmctl user create --local --email admin@example.com --username admin --password password123 --system-admin
```

**Populate with sample data** (optional):
```bash
bin/mmctl sampledata
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

**Important**: The web app isn't exposed directly - it's served via the server. Both server and web app must be running, and you access everything through `http://localhost:8065` (the server's port).

Open a **new terminal window** (keep the server running in the first).

```bash
# From the project root
cd webapp

# Install Node.js version (if using NVM - recommended)
nvm install

# Install npm dependencies (this takes a few minutes)
npm install

# Start the development server
make run
# OR alternatively: npm run dev-server
```

The web app will start and proxy API requests to the server running on port 8065.

**Expected Output**:
```
Webpack dev server starting...
Local:   http://localhost:8065/
```

**Note**: The web app runs on the same port as the server (8065) because it's proxied through the server.

---

## Step 6: Access the Web App

1. Open browser to `http://localhost:8065`
2. If you created an admin user via mmctl (Step 4), log in with those credentials
3. If not, you'll see the account creation screen:
   - Click **"Create an account"**
   - Fill in email, username, and password (must be 8+ characters)
   - Click **"Create Account"**
   - Create a team (e.g., "Dev Team")

**First user is automatically a System Admin** with full permissions.

**Alternative**: You can also add `http://localhost:8065` to the Mattermost desktop app.

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
# Webpack dev server auto-reloads automatically
# If you see issues, restart:
cd webapp
make run
# OR: npm run dev-server
```

### Stopping Services

**Stop the server**:
```bash
cd server
make stop-server
```

**Stop Docker containers** (separate command):
```bash
cd server
make stop-docker
```

**Note**: `stop-server` does NOT stop Docker containers. Use `stop-docker` for that.

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

### Customizing Behavior

You can customize server behavior by creating `server/config.override.mk`:
```bash
cd server
cp config.mk config.override.mk
# Edit config.override.mk to set options like MM_NO_DOCKER=true
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
make run               # Start dev server (recommended)
npm run dev-server     # Alternative: Start dev server
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
cd server && make run-server    # Starts server + Docker dependencies
cd webapp && make run            # Starts web app (in separate terminal)

# Stop everything
cd server
make stop-server                 # Stops server
make stop-docker                 # Stops Docker containers

# Reset everything
cd server
make stop-docker
docker compose down -v          # Remove volumes (deletes data)
rm -rf data                      # Remove local data directory
```

## Development Without Docker

If you prefer not to use Docker:

1. Copy `server/config.mk` to `server/config.override.mk`
2. Set `MM_NO_DOCKER=true` in `config.override.mk`
3. Install PostgreSQL locally and create database manually:
   ```bash
   psql postgres
   CREATE ROLE mmuser WITH LOGIN PASSWORD 'mostest';
   ALTER ROLE mmuser CREATEDB;
   \q
   psql postgres -U mmuser
   CREATE DATABASE mattermost_test;
   \q
   psql postgres
   GRANT ALL PRIVILEGES ON DATABASE mattermost_test TO mmuser;
   \q
   ```
4. Continue with normal setup steps

---

**You're all set!** 🚀

If you run into issues, check the [troubleshooting section](#common-issues-and-solutions) or ask in the Mattermost community forums.