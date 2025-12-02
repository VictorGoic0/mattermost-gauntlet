# Meeting Assistant Bot Plugin

AI-powered meeting transcription and summarization bot for Mattermost.

## Features

- Joins voice calls on-demand
- Records and transcribes audio using OpenAI Whisper
- Generates intelligent summaries with action items using GPT-4o
- Posts formatted results to channels or DMs

## Development

### Build Plugin Bundle

```bash
cd plugins-dev/meeting-bot
GOWORK=off make dist
```

This creates: `com.mattermost.meeting-bot-0.1.0.tar.gz`

### Upload to Mattermost

1. Go to System Console → Plugins → Plugin Management
2. Click "Upload Plugin"
3. Select the `.tar.gz` file
4. Click "Enable"

### Directory Structure

- **Source Code**: `plugins-dev/meeting-bot/` (edit here)
- **Extracted Plugin**: `server/plugins/com.mattermost.meeting-bot/` (managed by Mattermost, don't edit)

## Configuration

- **OpenAI API Key**: Required for transcription and summarization (coming soon)
- **Retention Days**: Default 30 days (user-configurable, coming soon)

## Usage

- `/meeting-bot start` - Start recording current call (coming soon)
- `/meeting-bot stop` - Stop recording (coming soon)
- `/meeting-bot settings` - Configure user settings (coming soon)

## Implementation Status

- ✅ Plugin skeleton
- ✅ Bot user creation
- ⏭️ Slash command registration (next)
- ⏭️ Database schema
- ⏭️ WebRTC recording
- ⏭️ Transcription
- ⏭️ Summarization

See `victor-features/1-meeting-bot/tasks-1.md` for detailed implementation tasks.

