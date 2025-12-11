# Context File Template

## PURPOSE OF THIS FILE
This file provides semantic context for working on a specific feature in the Mattermost codebase.
It is NOT a dump of entire files. It is a curated knowledge graph that helps Cursor:
1. Understand the relevant architecture
2. Know where code lives and why
3. Recognize patterns and conventions
4. Avoid common pitfalls

## HOW TO USE THIS FILE
When working on [FEATURE NAME], load this file into context.
This file should be the FIRST thing Cursor reads before any coding session.

---

## FEATURE OVERVIEW

We're building an AI Meeting Assistant Bot plugin that joins Mattermost voice calls, records audio, transcribes conversations using OpenAI Whisper, and generates intelligent summaries with extracted action items using GPT-4o. The bot operates on-demand (user-triggered via `/meeting-bot start`) and posts formatted results directly to the appropriate channel or DMs.

---

## ARCHITECTURE MAP

### System Components Involved
- **Mattermost Plugin System** - Our plugin extends Mattermost functionality
- **Calls Plugin** - External plugin that manages voice calls (uses rtcd service)
- **rtcd Service** - External WebRTC daemon that handles actual WebRTC connections
- **Pion WebRTC** - Go library we'll use to join calls independently (Option B approach)
- **OpenAI APIs** - Whisper (transcription) and GPT-4o (summarization)
- **Mattermost Database** - Store transcripts and summaries
- **Mattermost HTTP API** - Post messages, query call state

### Data Flow

**Recording Flow**:
```
User types /meeting-bot start
    ↓
Plugin validates: Is user in active call? (query Calls plugin or database)
    ↓
Plugin uses Pion WebRTC to join call as bot participant
    ↓
Plugin receives audio tracks, buffers audio data
    ↓
Call ends → Plugin saves audio to WAV file
    ↓
Plugin sends audio to Whisper API for transcription
    ↓
Plugin sends transcript to GPT-4o for summarization
    ↓
Plugin formats and posts result to channel/DM
    ↓
Plugin stores transcript in database
```

**Key Decision**: Using **Option B (Standalone WebRTC)** because Calls plugin doesn't expose APIs for external plugins to join calls programmatically.

---

## CODEBASE LOCATIONS

### Core Files (Must Understand)
- `plugins-dev/meeting-bot/plugin.go` - Main plugin file, OnActivate(), hooks
- `plugins-dev/meeting-bot/bot.go` - Bot user creation and management
- `plugins-dev/meeting-bot/commands.go` - Slash command handlers
- `plugins-dev/meeting-bot/recorder.go` - WebRTC join, audio recording
- `plugins-dev/meeting-bot/transcription.go` - Whisper API integration
- `plugins-dev/meeting-bot/summarization.go` - GPT-4o API integration
- `server/public/plugin/api.go` - Mattermost plugin API interface
- `server/channels/app/plugin_api.go` - Plugin HTTP communication (PluginHTTP method)

**Note**: Source code is in `plugins-dev/meeting-bot/` (root level) to avoid conflicts with Mattermost's plugin extraction directory `server/plugins/com.mattermost.meeting-bot/`.

### Calls Plugin Reference (research/calls-plugin/)
- `server/api.go` - HTTP API endpoints (can query `/calls/{channel_id}/active`)
- `server/rtcd.go` - RTCD service client (shows how Calls plugin connects to rtcd)
- `server/websocket.go` - WebSocket signaling (SDP, ICE candidates)
- `server/session.go` - Session management (how users join calls)

### Related Files (May Need to Touch)
- `server/channels/app/plugin_requests.go` - Inter-plugin HTTP communication
- `server/channels/db/migrations/` - Database migrations for our tables

### External Dependencies
- **Pion WebRTC** (https://github.com/pion/webrtc) - Pure Go WebRTC implementation for joining calls
- **OpenAI Whisper API** - Speech-to-text transcription
- **OpenAI GPT-4o API** - Summarization and action item extraction
- **go-audio/wav** - WAV file encoding for audio storage

---

## CODE PATTERNS & CONVENTIONS

### Pattern: Bot User Creation
**Where it's used**: `plugins-dev/meeting-bot/bot.go`
**What it does**: Creates or ensures a bot user exists when plugin activates. Uses `EnsureBotUser()` which is idempotent - safe to call multiple times.
**Example**:
```go
botID, err := p.API.EnsureBotUser(&model.Bot{
    Username:    "meeting-bot",
    DisplayName: "Meeting Assistant",
    Description: "AI-powered meeting transcription bot",
    OwnerId:     "com.mattermost.meeting-bot",
})
if err != nil {
    return "", fmt.Errorf("failed to ensure bot user: %w", err)
}
```
**Notes**: 
- `OwnerId` should be set to the plugin ID from manifest.json
- Store the returned `botID` in plugin struct for later use (posting messages, etc.)
- Call from `OnActivate()` hook

### Gotcha: [Name of Gotcha]
**What**: [What's the issue]
**Why**: [Why it happens]
**How to avoid**: [Solution]

---

## KEY FUNCTIONS & APIS

### Function: `API.EnsureBotUser()`
**Location**: `server/public/plugin/api.go` (plugin API interface)
**Purpose**: Creates a bot user if it doesn't exist, or returns existing bot user ID
**Parameters**: `*model.Bot` - Bot configuration (Username, DisplayName, Description, OwnerId)
**Returns**: `(string, error)` - Bot user ID and error
**When to use**: Call from `OnActivate()` to ensure bot user exists before plugin operations
**Example**:
```go
botID, err := p.API.EnsureBotUser(&model.Bot{
    Username:    "meeting-bot",
    DisplayName: "Meeting Assistant",
    Description: "AI-powered meeting transcription bot",
    OwnerId:     "com.mattermost.meeting-bot",
})
```

### Function: `FunctionName()`
**Location**: `path/to/file.go:123`
**Purpose**: [What it does]
**Parameters**: [What it takes]
**Returns**: [What it gives back]
**When to use**: [Use case]

---

## DATABASE SCHEMA

### Tables Involved
[Only tables relevant to this feature]

Example:
**Table**: `Posts`
- `Id` (string, primary key)
- `ChannelId` (string, foreign key to Channels)
- `Message` (text)
- `CreateAt` (bigint, unix timestamp)

---

## EXTERNAL DOCUMENTATION LINKS

### Official Mattermost Docs
- [Plugin Developer Guide](https://developers.mattermost.com/integrate/plugins/)
- [API Reference](https://api.mattermost.com/)

### Third-Party Docs
- [Pion WebRTC Examples](https://github.com/pion/webrtc/tree/master/examples)

---

## WORKING NOTES

### What I've Learned

**Calls Plugin Architecture**:
- Calls plugin uses external `rtcd` service for WebRTC (not handled directly in plugin)
- No public APIs for external plugins to join calls programmatically
- Bot APIs exist but only for Calls plugin's own bot (requires bot session auth)
- WebRTC signaling (SDP, ICE) goes through Mattermost WebSocket, then to rtcd service

**Plugin Communication**:
- Plugins can call other plugins via `PluginHTTP()` API
- Format: `/plugins/{plugin_id}/*`
- Can query Calls plugin endpoints like `/calls/{channel_id}/active`

**Decision: Option B (Standalone WebRTC)**:
- Cannot extend Calls plugin (no APIs exposed)
- Will use Pion WebRTC to join calls independently
- Need to figure out signaling mechanism (SDP/ICE exchange)
- May need to query database for call info or use HTTP API

**Bot User Creation** (2025-01-15):
- Use `p.API.EnsureBotUser(&model.Bot{...})` in `OnActivate()`
- `EnsureBotUser()` is idempotent - safe to call multiple times
- Must set `OwnerId` to plugin ID from manifest.json
- Store returned `botID` in plugin struct for later use (posting messages, etc.)
- Pattern: Create `bot.go` file with `ensureBotUser()` helper function

### Questions / Unknowns

1. **Signaling**: How do we exchange SDP offers/answers and ICE candidates?
   - Calls plugin uses rtcd service - do we connect to it or implement our own?
   
2. **TURN/STUN Servers**: How do we get ICE server configuration?
   - Calls plugin has `/turn-credentials` endpoint - is it accessible?

3. **Audio Format**: Confirm codec (likely Opus) and sample rate
   - Need to test actual audio format from WebRTC tracks

4. **Call Detection**: Best way to detect active calls?
   - HTTP API: `/plugins/com.mattermost.calls/calls/{channel_id}/active`
   - Database: Query `calls` table directly

5. **Call End Detection**: How do we know when call ends?
   - Poll database for call state?
   - Monitor WebSocket events (if accessible)?

### Debugging Tips
[To be added during implementation]

---

## LAST UPDATED
2025-01-15 - Added bot user creation pattern. Implemented `ensureBotUser()` in `bot.go`. Documented `API.EnsureBotUser()` usage pattern.