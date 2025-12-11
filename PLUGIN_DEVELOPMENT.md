# Plugin Development Setup

## Directory Structure

**Source Code (Development)**: `plugins-dev/meeting-bot/`
- This is where you edit code, build bundles, and develop
- **Never** gets touched by Mattermost
- Tracked in git
- Located at repo root to avoid conflicts

**Extracted Plugin (Runtime)**: `server/plugins/com.mattermost.meeting-bot/`
- This is where Mattermost extracts uploaded plugins
- **Gets deleted/replaced** when you upload a new version
- **Do NOT edit code here** - it will be overwritten
- Managed entirely by Mattermost

## Workflow

1. **Edit code** in `plugins-dev/meeting-bot/`
2. **Build bundle**: 
   ```bash
   cd plugins-dev/meeting-bot
   GOWORK=off make dist
   ```
3. **Upload bundle** via System Console → Plugins → Upload Plugin
4. Mattermost extracts to `server/plugins/com.mattermost.meeting-bot/`

## Why This Structure?

- **Prevents conflicts**: Mattermost manages `server/plugins/` and can delete/replace directories
- **Standard practice**: Plugins are typically developed in separate repositories
- **Git-friendly**: Source code stays safe and tracked
- **Clear separation**: Development vs. runtime directories are distinct

## Quick Reference

- **Build**: `cd plugins-dev/meeting-bot && GOWORK=off make dist`
- **Bundle location**: `plugins-dev/meeting-bot/com.mattermost.meeting-bot-0.1.0.tar.gz`
- **Upload**: System Console → Plugins → Plugin Management → Upload Plugin

