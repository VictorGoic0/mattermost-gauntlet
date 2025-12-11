# Meeting Assistant Bot Plugin

AI-powered meeting transcription and summarization bot for Mattermost.

## Overview

The Meeting Assistant Bot plugin joins Mattermost voice calls, records audio, transcribes conversations using OpenAI Whisper, and generates intelligent summaries with extracted action items using GPT-4o. The bot operates on-demand (user-triggered via `/meeting-bot start`) and posts formatted results directly to the appropriate channel.

## Prerequisites

- Mattermost Server 8.0.0 or higher
- PostgreSQL database (uses Mattermost's database)
- OpenAI API key (for transcription and summarization)
- Calls plugin enabled (for voice calls)

## Installation

### Build Plugin Bundle

```bash
cd plugins-dev/meeting-bot
GOWORK=off make dist
```

This creates: `com.mattermost.meeting-bot-0.1.0.tar.gz`

### Upload to Mattermost

1. Go to **System Console** → **Plugins** → **Plugin Management**
2. Click **"Upload Plugin"**
3. Select the `com.mattermost.meeting-bot-0.1.0.tar.gz` file
4. Click **"Enable"**

**Note**: Plugin uploads must be enabled in `config.json`:
```json
"PluginSettings": {
  "EnableUploads": true
}
```

### Directory Structure

- **Source Code**: `plugins-dev/meeting-bot/` (edit here, tracked in git)
- **Extracted Plugin**: `server/plugins/com.mattermost.meeting-bot/` (managed by Mattermost, don't edit)

## Configuration

Configure the plugin in **System Console** → **Plugins** → **Meeting Bot** → **Settings**:

- **OpenAI API Key** (required): Your OpenAI API key for Whisper and GPT-4o. Must start with `sk-`.
- **Enable Debug Logging** (optional): Enable additional debug logging for troubleshooting.

The API key is stored securely and masked in the UI.

## Usage

### Slash Commands

- `/meeting-bot help` - Show help message with available commands
- `/meeting-bot start` - Start recording the current call (coming in PR #2)
- `/meeting-bot stop` - Stop recording the current call (coming in PR #2)
- `/meeting-bot settings` - View or change your settings (coming in PR #7)

### Workflow

1. Join a voice call in a Mattermost channel
2. Type `/meeting-bot start` to begin recording
3. The bot joins the call and records audio
4. When the call ends, the bot processes the recording
5. Transcript and summary are posted to the channel

## Implementation Status

### ✅ Completed (PR #1)

- ✅ Plugin skeleton and structure
- ✅ Bot user creation
- ✅ Slash command registration
- ✅ Database schema (5 tables with indexes)
- ✅ Configuration management (OpenAI API key)

### ⏭️ Coming Soon

- ⏭️ WebRTC recording (PR #2)
- ⏭️ Audio transcription (PR #3)
- ⏭️ Summarization (PR #4)
- ⏭️ User settings (PR #7)

## Development

### Project Structure

```
plugins-dev/meeting-bot/
├── plugin.go          # Main plugin file, hooks
├── bot.go             # Bot user management
├── commands.go        # Slash command handlers
├── config.go          # Configuration management
├── storage.go         # Database initialization
├── plugin.json        # Plugin manifest
├── go.mod             # Go dependencies
├── Makefile           # Build system
└── README.md          # This file
```

### Building

```bash
cd plugins-dev/meeting-bot
GOWORK=off make dist
```

### Testing

After building, upload the plugin bundle via System Console. Check logs in:
- Docker Desktop (full logs)
- `server/logs/mattermost.log` (file logs)

## Database Schema

The plugin creates the following tables:

- `meeting_recordings` - Recording metadata
- `meeting_transcripts` - Transcript text
- `meeting_summaries` - Summaries with action items
- `plugin_config` - Plugin configuration storage
- `user_settings` - User preferences

All tables are created automatically on plugin activation.

## Troubleshooting

### Plugin won't activate
- Check server logs for errors
- Verify `EnableUploads` is true in config.json
- Check plugin bundle structure is correct

### Database tables not created
- Check logs for database initialization errors
- Verify PostgreSQL connection is working
- Check plugin has database access permissions

### Configuration not saving
- Verify you're a system admin
- Check API key format (should start with "sk-")
- Review logs for validation warnings

## See Also

- Detailed implementation tasks: `victor-features/1-meeting-bot/tasks-1.md`
- Investigation findings: `victor-features/1-meeting-bot/investigation.md`
- Feature context: `victor-context/features/1-meeting-bot/context.md`

