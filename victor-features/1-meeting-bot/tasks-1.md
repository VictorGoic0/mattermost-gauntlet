# tasks-1.md

## Overview
This file covers Investigation, Plugin Setup, and WebRTC Recording implementation.

**PRs Covered**: #1-2

---

## PR #1: Investigation & Plugin Setup

**Branch**: `feature/meeting-bot-setup`  
**Description**: Investigate Calls plugin architecture, set up plugin skeleton, create database schema, and register slash commands.

**Dependencies**: None (first PR)

**Success Criteria**:
- [ ] Calls plugin architecture documented (Option A vs B decided)
- [ ] Plugin loads successfully in Mattermost
- [ ] Bot user created and visible
- [ ] `/meeting-bot start` command registered and responds
- [ ] Database tables created
- [ ] All setup tests pass

---

## Investigation Phase

**Goal**: Understand Calls plugin architecture and decide integration approach.

**Reference**:
- 📄 TDD Section 3 (Technical Architecture)
- 📄 TDD Section 10 (Risks & Unknowns)

### Tasks:

- [ ] **1.1 Clone Calls plugin repository**
  - Repo: https://github.com/mattermost/mattermost-plugin-calls
  - Clone to separate directory for reference
  - Checkout latest stable version

- [ ] **1.2 Read Calls plugin documentation**
  - Review README and architecture docs
  - Understand how calls are initiated
  - Document call lifecycle (start → participants join → end)

- [ ] **1.3 Investigate WebRTC implementation**
  - Find WebRTC connection code in Calls plugin
  - Identify how participants join calls
  - Document WebRTC signaling process

- [ ] **1.4 Search for bot/recording APIs**
  - Search codebase for "bot", "record", "participant"
  - Check if API exists for programmatic join
  - Document any relevant API endpoints or hooks

- [ ] **1.5 Examine audio handling**
  - Find where audio streams are managed
  - Determine audio format (Opus, PCM, etc.)
  - Check if audio is pre-mixed or separate streams per participant

- [ ] **1.6 Check for call lifecycle events**
  - Search for event emissions: call_started, call_ended, etc.
  - Verify events are accessible to other plugins
  - Document event payloads and when they fire

- [ ] **1.7 Document findings**
  - Create investigation.md in features/1-meeting-bot/
  - Document Option A feasibility (extend Calls plugin)
  - Document Option B fallback (standalone WebRTC)
  - Make decision: Option A or Option B

- [ ] **1.8 If Option A: Document integration points**
  - List specific functions/APIs to use
  - Document required modifications to Calls plugin
  - Note any version compatibility issues

- [ ] **1.9 If Option B: Research Pion WebRTC**
  - Read Pion WebRTC examples: https://github.com/pion/webrtc/tree/master/examples
  - Find example for joining existing WebRTC session
  - Document signaling approach (how to get session details)

---

## Plugin Skeleton Setup

**Goal**: Create basic plugin structure that loads in Mattermost.

**Reference**:
- 📄 TDD Section 4 (Implementation Details - Plugin Structure)
- 📄 Context file: context/features/meeting-bot/context.md

### Tasks:

- [ ] **1.10 Create plugin directory structure**
  ```
  server/plugins/meeting-bot/
  ├── plugin.go
  ├── manifest.json
  └── README.md
  ```

- [ ] **1.11 Create manifest.json**
  - Set plugin ID: `com.mattermost.meeting-bot`
  - Set name: "Meeting Assistant Bot"
  - Set version: 0.1.0
  - Set min_server_version (check Mattermost version)
  - Set server executable path

- [ ] **1.12 Create plugin.go with basic structure**
  ```go
  package main
  
  import (
      "github.com/mattermost/mattermost-server/v6/plugin"
  )
  
  type Plugin struct {
      plugin.MattermostPlugin
      botUserID string
  }
  ```

- [ ] **1.13 Implement OnActivate() hook**
  - Log activation message
  - Initialize plugin state
  - Return nil (success) or error

- [ ] **1.14 Implement OnDeactivate() hook**
  - Log deactivation message
  - Clean up resources (if any at this stage)

- [ ] **1.15 Create Makefile for building plugin**
  - Add build target to compile Go code
  - Add deploy target to install locally
  - Add watch target for development

- [ ] **1.16 Build plugin binary**
  - Run `make build`
  - Verify binary created in dist/ folder
  - Check for compilation errors

- [ ] **1.17 Test plugin loads in Mattermost**
  - Upload plugin via System Console
  - Enable plugin
  - Check server logs for activation message
  - Verify no errors

---

## Bot User Creation

**Goal**: Create a system bot user for posting messages.

**Reference**:
- 📄 TDD Section 4 (Implementation Details - bot.go)

### Tasks:

- [ ] **1.18 Create bot.go file**
  - Add to `server/plugins/meeting-bot/bot.go`

- [ ] **1.19 Implement ensureBotUser() function**
  - Check if bot already exists (query by username)
  - If not exists, create bot user
  - Set username: "meeting-bot"
  - Set display name: "Meeting Assistant"
  - Set description: "AI-powered meeting transcription bot"

- [ ] **1.20 Set bot profile image (optional)**
  - Create bot avatar image (or use default)
  - Upload image via API
  - Set as bot profile picture

- [ ] **1.21 Call ensureBotUser() in OnActivate()**
  - Store bot user ID in plugin struct
  - Log bot user ID for debugging
  - Handle errors (bot creation failure)

- [ ] **1.22 Verify bot user visible**
  - Check Mattermost user list
  - Verify bot has "BOT" label
  - Verify bot cannot login (system user)

---

## Slash Command Registration

**Goal**: Register `/meeting-bot` slash command.

**Reference**:
- 📄 TDD Section 4 (Implementation Details - commands.go)

### Tasks:

- [ ] **1.23 Create commands.go file**
  - Add to `server/plugins/meeting-bot/commands.go`

- [ ] **1.24 Define command structure**
  ```go
  &model.Command{
      Trigger: "meeting-bot",
      AutoComplete: true,
      AutoCompleteDesc: "Control the meeting assistant bot",
      AutoCompleteHint: "[start|stop|settings|help]",
  }
  ```

- [ ] **1.25 Register command in OnActivate()**
  - Call `p.API.RegisterCommand()`
  - Handle registration errors
  - Log successful registration

- [ ] **1.26 Implement ExecuteCommand() hook**
  - Parse command arguments
  - Route to appropriate handler
  - Return command response

- [ ] **1.27 Implement handleHelp()**
  - Return help message with available commands
  - Format as ephemeral message (only user sees it)
  - Include usage examples

- [ ] **1.28 Implement handleStart() stub**
  - For now, just return "Start command received"
  - Will implement fully in PR #2
  - Return as ephemeral message

- [ ] **1.29 Implement handleStop() stub**
  - For now, just return "Stop command received"
  - Will implement fully in PR #2

- [ ] **1.30 Implement handleSettings() stub**
  - For now, just return "Settings command coming soon"
  - Will implement fully in PR #7

- [ ] **1.31 Test slash command**
  - Type `/meeting-bot` in any channel
  - Verify autocomplete appears
  - Run `/meeting-bot help`
  - Verify help message displays

- [ ] **1.32 Test command routing**
  - Run `/meeting-bot start`
  - Verify "Start command received" appears
  - Test other subcommands

---

## Database Schema Setup

**Goal**: Create database tables for storing recordings and transcripts.

**Reference**:
- 📄 TDD Section 5 (Database Schema)

### Tasks:

- [ ] **1.33 Create storage.go file**
  - Add to `server/plugins/meeting-bot/storage.go`

- [ ] **1.34 Implement initDatabase() function**
  - Check if tables already exist
  - If not, create all tables
  - Handle database errors

- [ ] **1.35 Create meeting_recordings table**
  ```sql
  CREATE TABLE IF NOT EXISTS meeting_recordings (
      id VARCHAR(26) PRIMARY KEY,
      call_id VARCHAR(26) NOT NULL,
      channel_id VARCHAR(26) NOT NULL,
      started_by VARCHAR(26) NOT NULL,
      started_at BIGINT NOT NULL,
      ended_at BIGINT,
      status VARCHAR(20) NOT NULL,
      audio_file_path TEXT,
      created_at BIGINT NOT NULL,
      updated_at BIGINT NOT NULL
  );
  ```

- [ ] **1.36 Add indexes to meeting_recordings**
  - Index on call_id
  - Index on status
  - Index on started_by

- [ ] **1.37 Create meeting_transcripts table**
  ```sql
  CREATE TABLE IF NOT EXISTS meeting_transcripts (
      id VARCHAR(26) PRIMARY KEY,
      call_id VARCHAR(26) NOT NULL,
      recording_id VARCHAR(26) NOT NULL,
      transcript TEXT NOT NULL,
      word_count INTEGER,
      created_at BIGINT NOT NULL,
      expires_at BIGINT
  );
  ```

- [ ] **1.38 Add indexes to meeting_transcripts**
  - Index on call_id
  - Index on expires_at (for cleanup job)
  - Index on created_at DESC (for queries)

- [ ] **1.39 Create meeting_summaries table**
  ```sql
  CREATE TABLE IF NOT EXISTS meeting_summaries (
      id VARCHAR(26) PRIMARY KEY,
      call_id VARCHAR(26) NOT NULL,
      recording_id VARCHAR(26) NOT NULL,
      transcript_id VARCHAR(26) NOT NULL,
      summary TEXT NOT NULL,
      action_items JSONB,
      topics JSONB,
      created_at BIGINT NOT NULL
  );
  ```

- [ ] **1.40 Add indexes to meeting_summaries**
  - Index on call_id
  - Index on recording_id

- [ ] **1.41 Create plugin_config table**
  ```sql
  CREATE TABLE IF NOT EXISTS plugin_config (
      key VARCHAR(255) PRIMARY KEY,
      value TEXT NOT NULL,
      encrypted BOOLEAN DEFAULT FALSE,
      updated_at BIGINT NOT NULL
  );
  ```

- [ ] **1.42 Create user_settings table**
  ```sql
  CREATE TABLE IF NOT EXISTS user_settings (
      user_id VARCHAR(26) PRIMARY KEY,
      retention_days INTEGER DEFAULT 30,
      auto_join BOOLEAN DEFAULT FALSE,
      created_at BIGINT NOT NULL,
      updated_at BIGINT NOT NULL
  );
  ```

- [ ] **1.43 Call initDatabase() in OnActivate()**
  - Run database initialization
  - Log table creation success
  - Handle errors gracefully

- [ ] **1.44 Test database creation**
  - Check PostgreSQL for created tables
  - Verify indexes exist
  - Test insert/query on each table

---

## Configuration Management

**Goal**: Set up API key configuration.

**Reference**:
- 📄 TDD Section 7 (Security & Permissions)

### Tasks:

- [ ] **1.45 Create config.go file**
  - Add to `server/plugins/meeting-bot/config.go`

- [ ] **1.46 Define configuration struct**
  ```go
  type Configuration struct {
      OpenAIAPIKey string
      EnableDebugLogging bool
  }
  ```

- [ ] **1.47 Implement OnConfigurationChange() hook**
  - Called when admin updates settings
  - Reload configuration
  - Validate API key format

- [ ] **1.48 Implement getOpenAIKey() helper**
  - Retrieve API key from configuration
  - Return error if not configured
  - Never log the actual key

- [ ] **1.49 Add configuration validation**
  - Check if API key is set
  - Validate key format (starts with "sk-")
  - Warn admin if key missing

- [ ] **1.50 Create plugin settings schema (manifest.json)**
  ```json
  "settings_schema": {
      "settings": [
          {
              "key": "OpenAIAPIKey",
              "type": "text",
              "display_name": "OpenAI API Key",
              "help_text": "Enter your OpenAI API key"
          }
      ]
  }
  ```

- [ ] **1.51 Test configuration in System Console**
  - Navigate to Plugins > Meeting Bot > Settings
  - Enter test API key
  - Save and verify key stored

---

## PR #1 Wrap-Up

- [ ] **1.52 Create README.md for plugin**
  - Document what plugin does
  - List prerequisites
  - Add installation instructions
  - Include configuration steps

- [ ] **1.53 Add basic logging**
  - Log plugin activation
  - Log bot user creation
  - Log command registration
  - Log database initialization

- [ ] **1.54 Test full plugin lifecycle**
  - Upload and enable plugin
  - Verify all initialization steps succeed
  - Run `/meeting-bot help`
  - Disable and re-enable plugin

- [ ] **1.55 Commit and create PR**
  - Git add all files
  - Commit: "feat: Add meeting bot plugin skeleton and database setup"
  - Push to branch
  - Create PR with description

---

**PR #1 Complete** ✅

---

## PR #2: WebRTC Integration & Recording

**Branch**: `feature/meeting-bot-recording`  
**Description**: Implement WebRTC call joining, audio capture, and recording functionality.

**Dependencies**: PR #1 merged

**Success Criteria**:
- [ ] Bot can join active calls
- [ ] Audio is captured from all participants
- [ ] Audio is saved to WAV file on call end
- [ ] Multiple concurrent recordings supported
- [ ] No audio data loss or corruption
- [ ] Integration tests pass

---

## WebRTC Setup

**Goal**: Set up Pion WebRTC library and basic connection.

**Reference**:
- 📄 TDD Section 3 (Technical Architecture - WebRTC)
- 📄 Investigation findings from PR #1

### Tasks:

- [ ] **2.1 Add Pion WebRTC dependency**
  - Run `go get github.com/pion/webrtc/v3`
  - Update go.mod and go.sum
  - Verify import works

- [ ] **2.2 Create recorder.go file**
  - Add to `server/plugins/meeting-bot/recorder.go`

- [ ] **2.3 Define Recording struct**
  ```go
  type Recording struct {
      CallID         string
      ChannelID      string
      StartedBy      string
      StartTime      time.Time
      PeerConnection *webrtc.PeerConnection
      AudioBuffer    *bytes.Buffer
      Status         string
      mu             sync.Mutex
  }
  ```

- [ ] **2.4 Add active recordings map to Plugin struct**
  ```go
  activeRecordings map[string]*Recording
  recordingsMu     sync.RWMutex
  ```

- [ ] **2.5 Initialize recordings map in OnActivate()**
  - Create empty map
  - Initialize mutex

---

## Calls Plugin Integration

**Goal**: Integrate with Calls plugin to get call information.

**Reference**: 
- 📄 Investigation findings (Option A or B)
- 📄 Calls plugin documentation

### Tasks:

- [ ] **2.6 Implement getActiveCallInChannel()**
  - Query Calls plugin for active call in channel
  - Return call details (call_id, participants, WebRTC info)
  - Handle "no active call" case

- [ ] **2.7 If Option A: Use Calls plugin API**
  - Call Calls plugin's exposed API
  - Get WebRTC signaling information
  - Get ICE servers configuration

- [ ] **2.8 If Option B: Implement standalone approach**
  - Query Mattermost for call metadata
  - Construct WebRTC offer independently
  - Handle signaling manually

- [ ] **2.9 Implement getCallParticipants()**
  - Get list of users in call
  - Return user IDs and usernames
  - Handle empty call case

- [ ] **2.10 Subscribe to call events**
  - Listen for call_started event
  - Listen for call_ended event
  - Listen for user_joined_call event
  - Listen for user_left_call event

- [ ] **2.11 Implement call_ended event handler**
  - When call ends, trigger processing
  - Get recording for that call_id
  - Call processRecording() function

---

## Bot Join Implementation

**Goal**: Bot joins WebRTC call as participant.

**Reference**: 
- 📄 TDD Section 4 (Implementation Details - recorder.go JoinCall)
- 📄 Pion WebRTC examples

### Tasks:

- [ ] **2.12 Implement JoinCall() function signature**
  ```go
  func (p *Plugin) JoinCall(callID string, callInfo CallDetails) error
  ```

- [ ] **2.13 Create WebRTC peer connection**
  - Configure ICE servers (from call info)
  - Set up peer connection config
  - Create peer connection object

- [ ] **2.14 Set up OnTrack handler**
  - Register handler for incoming audio tracks
  - Pass to handleAudioTrack() function
  - Log track details

- [ ] **2.15 Set up OnICECandidate handler**
  - Send ICE candidates to Calls plugin
  - Handle ICE gathering state changes
  - Log ICE connection state

- [ ] **2.16 Create and send WebRTC offer**
  - Generate SDP offer
  - Set as local description
  - Send to Calls plugin signaling server

- [ ] **2.17 Receive and process answer**
  - Get SDP answer from Calls plugin
  - Set as remote description
  - Wait for ICE connection

- [ ] **2.18 Create Recording object**
  - Initialize with call details
  - Set status to "recording"
  - Add to activeRecordings map

- [ ] **2.19 Insert recording to database**
  - Save to meeting_recordings table
  - Set status = "recording"
  - Store call_id, channel_id, started_by, started_at

- [ ] **2.20 Return success or error**
  - Return nil on success
  - Return descriptive error on failure

---

## Audio Capture Implementation

**Goal**: Capture audio packets and buffer them.

**Reference**: 
- 📄 TDD Section 4 (Implementation Details - handleAudioTrack)

### Tasks:

- [ ] **2.21 Implement handleAudioTrack() function**
  ```go
  func (p *Plugin) handleAudioTrack(callID string, track *webrtc.TrackRemote)
  ```

- [ ] **2.22 Get recording from activeRecordings**
  - Lock map with mutex
  - Retrieve Recording by callID
  - Unlock mutex

- [ ] **2.23 Read RTP packets from track**
  - Loop: `track.ReadRTP()`
  - Handle EOF (track closed)
  - Continue reading until track ends

- [ ] **2.24 Decode Opus to PCM**
  - Initialize Opus decoder
  - Decode RTP payload to PCM samples
  - Handle decode errors

- [ ] **2.25 Write PCM to audio buffer**
  - Lock Recording mutex
  - Append PCM data to AudioBuffer
  - Unlock mutex

- [ ] **2.26 Handle multiple audio tracks**
  - Each participant = separate track
  - Mix tracks together (simple addition)
  - Normalize to prevent clipping

- [ ] **2.27 Handle track errors**
  - Log track read errors
  - Continue with other tracks
  - Don't crash on single track failure

---

## Audio Processing

**Goal**: Convert audio buffer to WAV file.

**Reference**: 
- 📄 TDD Section 4 (Implementation Details - SaveRecording)

### Tasks:

- [ ] **2.28 Add audio library dependencies**
  - `go get github.com/go-audio/audio`
  - `go get github.com/go-audio/wav`

- [ ] **2.29 Implement SaveRecording() function**
  ```go
  func (p *Plugin) SaveRecording(callID string) (string, error)
  ```

- [ ] **2.30 Get recording from activeRecordings**
  - Lock and retrieve Recording
  - Check recording exists
  - Get AudioBuffer

- [ ] **2.31 Create temp WAV file**
  - Path: `/tmp/meeting-{callID}.wav`
  - Create file with write permissions
  - Defer file close

- [ ] **2.32 Write WAV header**
  - Sample rate: 16000 Hz (Whisper requirement)
  - Bit depth: 16-bit
  - Channels: 1 (mono)
  - Create wav.Encoder

- [ ] **2.33 Convert buffer to audio.IntBuffer**
  - Read PCM data from AudioBuffer
  - Create audio.Format struct
  - Wrap in audio.IntBuffer

- [ ] **2.34 Write audio data to WAV**
  - Call encoder.Write()
  - Handle write errors
  - Close encoder

- [ ] **2.35 Verify WAV file created**
  - Check file exists
  - Check file size > 0
  - Log file path and size

- [ ] **2.36 Return file path**
  - Return `/tmp/meeting-{callID}.wav`
  - Return error if failed

---

## Command Handler Updates

**Goal**: Implement `/meeting-bot start` command fully.

**Reference**: 
- 📄 TDD Section 2 (User Experience - Starting Recording)
- 📄 TDD Section 4 (commands.go - handleStartRecording)

### Tasks:

- [ ] **2.37 Update handleStart() in commands.go**
  - Remove stub implementation
  - Add full logic

- [ ] **2.38 Validate user is in active call**
  - Call getActiveCallInChannel()
  - If no call, return error message
  - Message: "❌ No active call detected in this channel. Please join a call first."

- [ ] **2.39 Check if already recording**
  - Check activeRecordings map for call_id
  - If exists, return message
  - Message: "🤖 Already recording this call (started by @username at HH:MM)"

- [ ] **2.40 Start recording**
  - Call JoinCall() function
  - Handle join errors
  - If error, post error message

- [ ] **2.41 Post confirmation message**
  - Post to channel (not ephemeral)
  - Message: "🤖 Meeting Assistant joined the call and is recording"
  - Use bot user ID as author

- [ ] **2.42 Handle mid-call join notification**
  - If call already in progress, note join time
  - Message: "🤖 Recording started. Note: I joined at HH:MM, so earlier conversation is not recorded."

---

## Recording Lifecycle Management

**Goal**: Handle recording from start to end.

**Reference**: 
- 📄 TDD Section 3 (Data Flow - Recording Flow)

### Tasks:

- [ ] **2.43 Implement stopRecording() function**
  - Get recording from map
  - Close WebRTC peer connection
  - Set status to "processing"

- [ ] **2.44 Handle call_ended event**
  - When event fires, get call_id
  - Find recording in activeRecordings
  - Call stopRecording()

- [ ] **2.45 Save audio on call end**
  - Call SaveRecording()
  - Get file path
  - Update recording in database (set ended_at, audio_file_path)

- [ ] **2.46 Post "Processing..." message**
  - Post to channel
  - Message: "🤖 Processing your meeting... This may take 2 minutes."

- [ ] **2.47 Remove from activeRecordings**
  - Delete from map after processing starts
  - Keep in database with status "processing"

- [ ] **2.48 Handle unexpected disconnections**
  - If WebRTC connection drops, log warning
  - Save whatever audio was captured
  - Post message about incomplete recording

---

## Multiple Recordings Support

**Goal**: Support multiple concurrent recordings.

**Reference**: 
- 📄 TDD Section 10 (Risks - Multiple Concurrent Calls)

### Tasks:

- [ ] **2.49 Test concurrent recordings**
  - Start recording in channel A
  - Start recording in channel B simultaneously
  - Verify both record independently

- [ ] **2.50 Add thread-safe map access**
  - Use RWMutex for activeRecordings
  - Lock for writes, RLock for reads
  - Prevent race conditions

- [ ] **2.51 Ensure separate audio buffers**
  - Each Recording has own AudioBuffer
  - No shared state between recordings
  - Verify no audio mixing between calls

- [ ] **2.52 Test recording cleanup**
  - Start 3 recordings
  - End them in different order
  - Verify all cleaned up properly

---

## Error Handling

**Goal**: Handle WebRTC and recording errors gracefully.

**Reference**: 
- 📄 TDD Section 9 (Error Handling)

### Tasks:

- [ ] **2.53 Handle "no active call" error**
  - Return user-friendly message
  - Don't crash plugin

- [ ] **2.54 Handle WebRTC connection failure**
  - Log detailed error
  - Post message: "❌ Failed to join call. Please try again."
  - Clean up partial recording

- [ ] **2.55 Handle audio track errors**
  - Continue recording other tracks
  - Log which track failed
  - Don't stop entire recording

- [ ] **2.56 Handle disk full error**
  - Catch error when saving WAV
  - Post message about disk space
  - Alert admin

- [ ] **2.57 Handle call ending before bot joins**
  - Check if call still active after join attempt
  - If ended, clean up and notify user

---

## Testing

**Goal**: Verify recording functionality works.

### Tasks:

- [ ] **2.58 Manual test: Join and record call**
  - Start voice call in test channel
  - Run `/meeting-bot start`
  - Verify bot joins
  - Talk for 2 minutes
  - End call

- [ ] **2.59 Verify audio file created**
  - Check `/tmp/` for WAV file
  - Open file in audio player
  - Verify audio is clear and audible

- [ ] **2.60 Test "already recording" message**
  - Start recording
  - Try to start again
  - Verify error message

- [ ] **2.61 Test "no active call" message**
  - Run command without active call
  - Verify error message

- [ ] **2.62 Test concurrent recordings**
  - Start 2 calls in different channels
  - Record both simultaneously
  - Verify separate audio files

- [ ] **2.63 Test mid-call join**
  - Start call, talk for 1 minute
  - Start recording
  - Verify notification about join time

---

## PR #2 Wrap-Up

- [ ] **2.64 Add detailed logging**
  - Log WebRTC connection state
  - Log audio track details
  - Log file save success/failure

- [ ] **2.65 Update README**
  - Document recording functionality
  - Add troubleshooting section

- [ ] **2.66 Commit and create PR**
  - Git add all files
  - Commit: "feat: Implement WebRTC recording functionality"
  - Push to branch
  - Create PR with description and demo video

---

**PR #2 Complete** ✅
