# AI Meeting Assistant Bot - Technical Design Document

## Metadata
- **Priority**: 1 (Highest)
- **Estimated Timeline**: 7-10 days
- **Status**: Planning
- **Last Updated**: 2025-01-15

---

## 1. Feature Overview

### What We're Building
An AI-powered meeting assistant bot that joins Mattermost voice calls, records audio, transcribes conversations using OpenAI Whisper, and generates intelligent summaries with extracted action items using GPT-4o. The bot operates on-demand (user-triggered) and posts formatted results directly to the appropriate channel or DMs.

### Why This Matters
- **Time savings**: Eliminates manual note-taking during meetings (30+ minutes per meeting)
- **Accountability**: Automatically extracts and tracks action items
- **Accessibility**: Provides written record of verbal discussions
- **Knowledge capture**: Meeting decisions and context preserved searchable
- **Universal need**: Every team has meetings, everyone forgets details

### Success Criteria
- [ ] Bot successfully joins active calls via user command
- [ ] Bot records audio from all participants (mixed stream)
- [ ] Transcription accuracy >90% for clear audio
- [ ] Action items correctly extracted from transcript
- [ ] Summary posted within 2 minutes of call ending
- [ ] Supports multiple concurrent recordings
- [ ] Zero data loss (no failed transcriptions without retry)
- [ ] User-configurable data retention settings

---

## 2. User Experience

### User Flow

#### Starting a Recording

**Scenario**: Team having a voice call in #engineering channel

1. **User** (in active call): Types `/meeting-bot start` in channel
2. **Bot**: Posts message "🤖 Meeting Assistant joined the call and is recording"
3. **User**: Continues call normally (bot is invisible/silent participant)
4. **Call continues**: Bot records audio in background

#### Call Ends

5. **Call ends**: Mattermost Calls plugin fires `call_ended` event
6. **Bot**: Immediately posts "🤖 Processing your meeting... This may take 2 minutes."
7. **Bot**: 
   - Saves audio to temp storage (`/tmp/meeting-{id}.wav`)
   - Sends audio to Whisper API for transcription
   - If API fails: Posts "⚠️ Transcription failed, retrying... (Attempt 1/3)"
   - Receives transcript
   - Deletes audio file immediately
8. **Bot**:
   - Sends transcript to GPT-4o with prompt to extract action items + summary
   - Formats result
9. **Bot**: Posts final formatted message (see format below)

#### Result Format

```markdown
## 🎯 Meeting Summary

**Duration**: 32 minutes  
**Participants**: @alice, @bob, @charlie  
**Started**: 2:15 PM  
**Ended**: 2:47 PM

---

### ✅ Action Items

- @alice to send proposal by Friday
- @bob to review budget numbers
- @charlie to schedule follow-up meeting next week
- Research competitor pricing (unassigned)

---

### 📝 Summary

We discussed Q4 budget allocation and decided to increase marketing spend by 15%. The team agreed to pilot the new CRM system starting next month. Key concerns were raised about staffing for the project. Bob will review final numbers before Friday's deadline.

---

### 📄 Full Transcript

<details>
<summary>Click to expand full transcript (32 minutes)</summary>

Hey everyone, let's get started with the Q4 budget discussion. I've been looking at the numbers and I think we need to increase marketing spend by about 15%. The current allocation isn't enough to hit our growth targets.

That makes sense. I can review the budget numbers and make sure we have room for that increase. What about the CRM system we've been discussing?

Good point. I think we're ready to pilot it starting next month. We should probably schedule a follow-up meeting next week to nail down the timeline.

[Transcript continues...]

</details>
```

#### Result Delivery Logic

**If call was in a channel** (e.g., #engineering):
- Post result in that channel

**If call was a group DM** (e.g., Alice, Bob, Charlie):
- Post result in that group DM

**If call was 1-on-1 DM**:
- Post result in that DM

### UI/UX Requirements

**MVP (Slash Command)**:
- User types `/meeting-bot start` to trigger
- Bot responds with confirmation message
- Bot posts status updates during processing
- Final result posted as rich markdown message

**Post-MVP (Button in Call UI)**:
- "Assistant!" button visible in active call window
- Clicking button triggers same flow as slash command
- Requires modifying Calls plugin UI (separate repo)

### Edge Cases & User Feedback

**If bot joins mid-call**:
- Bot posts: "🤖 Recording started. Note: I joined at 2:10 PM, so earlier conversation is not recorded."

**If API fails after 3 retries**:
- Bot posts: "❌ Transcription failed after 3 attempts. Audio has been saved. Please contact admin or try `/meeting-bot retry {recording_id}`"

**If no speech detected**:
- Bot posts: "⚠️ No speech detected in recording. This might be a technical issue. Please verify your audio setup."

**If multiple people try to start bot in same call**:
- Bot posts: "🤖 Already recording this call (started by @alice at 2:10 PM)"

**If user tries to start bot when not in a call**:
- Bot posts: "❌ No active call detected in this channel. Please join a call first."

---

## 3. Technical Architecture

### System Components Involved

```
┌─────────────────────────────────────────────────────────────┐
│                   MATTERMOST CALLS PLUGIN                   │
│  - Manages WebRTC connections                               │
│  - Handles call lifecycle (start, end)                      │
│  - Fires events: call_started, call_ended                   │
│  - Provides audio streams from participants                 │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  │ Events & Audio Streams
                  ▼
┌─────────────────────────────────────────────────────────────┐
│              MEETING BOT PLUGIN (OUR CODE)                  │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐  │
│  │  Command Handler                                     │  │
│  │  - Listens for /meeting-bot start                   │  │
│  │  - Validates user in active call                    │  │
│  │  - Triggers bot join logic                          │  │
│  └─────────────────────────────────────────────────────┘  │
│                          │                                  │
│                          ▼                                  │
│  ┌─────────────────────────────────────────────────────┐  │
│  │  Call Recorder                                       │  │
│  │  - Joins WebRTC session as participant              │  │
│  │  - Receives audio streams (all participants)        │  │
│  │  - Mixes audio into single file                     │  │
│  │  - Saves to /tmp/meeting-{id}.wav                   │  │
│  └─────────────────────────────────────────────────────┘  │
│                          │                                  │
│                          ▼                                  │
│  ┌─────────────────────────────────────────────────────┐  │
│  │  Transcription Service                               │  │
│  │  - Sends audio to OpenAI Whisper API                │  │
│  │  - Handles retries (3 attempts)                     │  │
│  │  - Returns transcript text                          │  │
│  │  - Deletes audio file immediately                   │  │
│  └─────────────────────────────────────────────────────┘  │
│                          │                                  │
│                          ▼                                  │
│  ┌─────────────────────────────────────────────────────┐  │
│  │  Summarization Service                               │  │
│  │  - Sends transcript to GPT-4o                        │  │
│  │  - Extracts: action items, summary, topics          │  │
│  │  - Returns structured JSON                          │  │
│  └─────────────────────────────────────────────────────┘  │
│                          │                                  │
│                          ▼                                  │
│  ┌─────────────────────────────────────────────────────┐  │
│  │  Result Formatter & Poster                           │  │
│  │  - Formats markdown message                          │  │
│  │  - Determines delivery location (channel/DM)        │  │
│  │  - Posts via Mattermost API                         │  │
│  │  - Stores transcript in database                    │  │
│  └─────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                  EXTERNAL SERVICES                          │
│                                                             │
│  ┌──────────────────┐         ┌──────────────────┐        │
│  │  OpenAI Whisper  │         │     GPT-4o       │        │
│  │  API             │         │     API          │        │
│  │  (Transcription) │         │  (Summarization) │        │
│  └──────────────────┘         └──────────────────┘        │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                   POSTGRESQL DATABASE                       │
│  - meeting_recordings table                                 │
│  - meeting_transcripts table                                │
│  - meeting_summaries table                                  │
│  - plugin_config table (API keys, settings)                 │
└─────────────────────────────────────────────────────────────┘
```

### Data Flow

**1. Bot Join Flow**:
```
User types /meeting-bot start
    ↓
Plugin validates: Is user in active call?
    ↓
Plugin queries Calls plugin: Get call details (call_id, participants, WebRTC info)
    ↓
Plugin joins WebRTC session as "bot participant"
    ↓
Plugin posts confirmation message
    ↓
Plugin starts recording audio stream
```

**2. Recording Flow**:
```
Audio streams from all participants
    ↓
Mix into single mono/stereo stream
    ↓
Write to buffer in memory
    ↓
On call_ended event: Flush buffer to /tmp/meeting-{id}.wav
```

**3. Processing Flow**:
```
Call ends → call_ended event fires
    ↓
Plugin posts "Processing..." message
    ↓
Save audio to /tmp/meeting-{id}.wav
    ↓
Send audio file to Whisper API (POST multipart/form-data)
    ↓
Whisper returns transcript JSON: { "text": "..." }
    ↓
Delete audio file from /tmp/
    ↓
Send transcript to GPT-4o with structured prompt
    ↓
GPT-4o returns JSON: { "action_items": [...], "summary": "...", "topics": [...] }
    ↓
Format markdown message
    ↓
Determine delivery location (channel vs DM)
    ↓
Post message via Mattermost API
    ↓
Store transcript in database (with user's retention setting)
```

### External Dependencies

**1. Mattermost Calls Plugin**
- **What**: Official plugin for voice/video calls
- **Repo**: https://github.com/mattermost/mattermost-plugin-calls
- **Why we need it**: Provides WebRTC infrastructure for calls
- **Version**: Latest (check compatibility)
- **Integration approach**: 
  - Option A (preferred): Extend Calls plugin to expose bot join API
  - Option B (fallback): Join WebRTC independently using Pion

**2. Pion WebRTC (Go library)**
- **What**: Pure Go implementation of WebRTC
- **Package**: `github.com/pion/webrtc/v3`
- **Why we need it**: To join calls and receive audio streams
- **Version**: v3.x (latest stable)

**3. OpenAI Whisper API**
- **What**: Speech-to-text transcription service
- **Endpoint**: `https://api.openai.com/v1/audio/transcriptions`
- **Why we need it**: Convert audio to text
- **Cost**: $0.006 per minute of audio
- **API Key**: Provided (you have access)
- **Rate limits**: 50 requests/minute

**4. OpenAI GPT-4o API**
- **What**: LLM for summarization and extraction
- **Endpoint**: `https://api.openai.com/v1/chat/completions`
- **Model**: `gpt-4o`
- **Why we need it**: Extract action items and generate summary
- **Cost**: $0.0025 per 1K tokens (input/output)
- **API Key**: Same as Whisper (OpenAI unified key)

**5. Go Libraries**
```go
// Core plugin
github.com/mattermost/mattermost-server/v6/plugin
github.com/mattermost/mattermost-server/v6/model

// WebRTC
github.com/pion/webrtc/v3

// HTTP client
github.com/go-resty/resty/v2

// Audio processing
github.com/go-audio/audio
github.com/go-audio/wav
```

---

## 4. Implementation Details

### Plugin Structure

```
server/plugins/meeting-bot/
├── plugin.go              # Main plugin file, OnActivate(), hooks
├── manifest.json          # Plugin metadata
├── bot.go                 # Bot user creation and management
├── commands.go            # Slash command handlers
├── recorder.go            # WebRTC join, audio recording
├── transcription.go       # Whisper API integration
├── summarization.go       # GPT-4o API integration
├── formatter.go           # Markdown formatting
├── storage.go             # Database operations
├── config.go              # Configuration management
├── utils.go               # Helper functions
└── README.md              # Plugin documentation

webapp/plugins/meeting-bot/  (Post-MVP for UI button)
└── index.tsx              # React component for "Assistant!" button
```

### Key Functions & APIs

#### 1. Plugin Lifecycle

**Function**: `OnActivate()`
- **Location**: `plugin.go`
- **Purpose**: Initialize plugin when activated
- **Implementation**:
  ```go
  func (p *Plugin) OnActivate() error {
      // Register slash commands
      p.API.RegisterCommand(&model.Command{
          Trigger: "meeting-bot",
          AutoComplete: true,
          AutoCompleteDesc: "Control the meeting assistant bot",
      })
      
      // Create bot user if doesn't exist
      botUser, err := p.ensureBotUser()
      if err != nil {
          return err
      }
      p.botUserID = botUser.Id
      
      // Subscribe to Calls plugin events
      p.subscribeToCallEvents()
      
      // Initialize database tables
      p.initDatabase()
      
      return nil
  }
  ```

#### 2. Command Handler

**Function**: `ExecuteCommand()`
- **Location**: `commands.go`
- **Purpose**: Handle `/meeting-bot start` command
- **Parameters**: 
  - `args *model.CommandArgs` (contains channel_id, user_id, command text)
- **Returns**: `(*model.CommandResponse, *model.AppError)`
- **Implementation**:
  ```go
  func (p *Plugin) ExecuteCommand(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
      command := args.Command
      
      if strings.HasPrefix(command, "/meeting-bot start") {
          return p.handleStartRecording(args)
      }
      
      if strings.HasPrefix(command, "/meeting-bot stop") {
          return p.handleStopRecording(args)
      }
      
      if strings.HasPrefix(command, "/meeting-bot settings") {
          return p.handleSettings(args)
      }
      
      return p.showHelp(args), nil
  }
  
  func (p *Plugin) handleStartRecording(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
      // 1. Check if user is in active call
      callInfo, err := p.getActiveCallInChannel(args.ChannelId)
      if err != nil {
          return p.ephemeralResponse("No active call detected in this channel."), nil
      }
      
      // 2. Check if already recording
      if p.isRecording(callInfo.ID) {
          return p.ephemeralResponse("Already recording this call."), nil
      }
      
      // 3. Start recording
      err = p.startRecording(callInfo, args.UserId)
      if err != nil {
          return p.ephemeralResponse("Failed to start recording: " + err.Error()), nil
      }
      
      // 4. Post confirmation
      p.postMessage(args.ChannelId, "🤖 Meeting Assistant joined the call and is recording")
      
      return &model.CommandResponse{}, nil
  }
  ```

#### 3. WebRTC Integration

**Function**: `JoinCall()`
- **Location**: `recorder.go`
- **Purpose**: Join WebRTC call as bot participant
- **Parameters**:
  - `callID string`
  - `callInfo CallDetails` (WebRTC signaling info)
- **Returns**: `error`
- **Implementation**:
  ```go
  func (p *Plugin) JoinCall(callID string, callInfo CallDetails) error {
      // Create WebRTC peer connection
      peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
          ICEServers: callInfo.ICEServers,
      })
      if err != nil {
          return err
      }
      
      // Handle incoming audio tracks
      peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
          p.handleAudioTrack(callID, track)
      })
      
      // Create offer
      offer, err := peerConnection.CreateOffer(nil)
      if err != nil {
          return err
      }
      
      // Set local description
      err = peerConnection.SetLocalDescription(offer)
      if err != nil {
          return err
      }
      
      // Send offer to Calls plugin signaling server
      answer, err := p.sendOfferToCallsPlugin(callID, offer)
      if err != nil {
          return err
      }
      
      // Set remote description
      err = peerConnection.SetRemoteDescription(answer)
      if err != nil {
          return err
      }
      
      // Store connection
      p.activeRecordings[callID] = &Recording{
          CallID: callID,
          PeerConnection: peerConnection,
          AudioBuffer: &bytes.Buffer{},
          StartTime: time.Now(),
      }
      
      return nil
  }
  ```

#### 4. Audio Recording

**Function**: `handleAudioTrack()`
- **Location**: `recorder.go`
- **Purpose**: Receive audio packets and write to buffer
- **Implementation**:
  ```go
  func (p *Plugin) handleAudioTrack(callID string, track *webrtc.TrackRemote) {
      recording := p.activeRecordings[callID]
      
      // Read RTP packets from track
      for {
          rtpPacket, _, err := track.ReadRTP()
          if err != nil {
              break
          }
          
          // Decode audio (Opus -> PCM)
          pcmData, err := p.decodeOpus(rtpPacket.Payload)
          if err != nil {
              continue
          }
          
          // Write to buffer
          recording.AudioBuffer.Write(pcmData)
      }
  }
  ```

**Function**: `SaveRecording()`
- **Location**: `recorder.go`
- **Purpose**: Save audio buffer to WAV file
- **Implementation**:
  ```go
  func (p *Plugin) SaveRecording(callID string) (string, error) {
      recording := p.activeRecordings[callID]
      
      // Create WAV file
      filename := fmt.Sprintf("/tmp/meeting-%s.wav", callID)
      file, err := os.Create(filename)
      if err != nil {
          return "", err
      }
      defer file.Close()
      
      // Write WAV header
      encoder := wav.NewEncoder(file, 16000, 16, 1, 1) // 16kHz, 16-bit, mono
      
      // Write audio data
      buf := &audio.IntBuffer{
          Data: recording.AudioBuffer.Bytes(),
          Format: &audio.Format{
              SampleRate: 16000,
              NumChannels: 1,
          },
      }
      
      err = encoder.Write(buf)
      if err != nil {
          return "", err
      }
      
      err = encoder.Close()
      if err != nil {
          return "", err
      }
      
      return filename, nil
  }
  ```

#### 5. Transcription

**Function**: `TranscribeAudio()`
- **Location**: `transcription.go`
- **Purpose**: Send audio to Whisper API and get transcript
- **Parameters**:
  - `audioPath string` (path to WAV file)
- **Returns**: `(string, error)` (transcript text, error)
- **Implementation**:
  ```go
  func (p *Plugin) TranscribeAudio(audioPath string) (string, error) {
      // Read audio file
      audioData, err := ioutil.ReadFile(audioPath)
      if err != nil {
          return "", err
      }
      
      // Prepare request
      client := resty.New()
      resp, err := client.R().
          SetHeader("Authorization", "Bearer "+p.getOpenAIKey()).
          SetFileReader("file", filepath.Base(audioPath), bytes.NewReader(audioData)).
          SetFormData(map[string]string{
              "model": "whisper-1",
              "response_format": "json",
          }).
          Post("https://api.openai.com/v1/audio/transcriptions")
      
      if err != nil {
          return "", err
      }
      
      if resp.StatusCode() != 200 {
          return "", fmt.Errorf("Whisper API error: %s", resp.String())
      }
      
      // Parse response
      var result struct {
          Text string `json:"text"`
      }
      
      err = json.Unmarshal(resp.Body(), &result)
      if err != nil {
          return "", err
      }
      
      return result.Text, nil
  }
  ```

**Function**: `TranscribeWithRetry()`
- **Location**: `transcription.go`
- **Purpose**: Retry transcription up to 3 times
- **Implementation**:
  ```go
  func (p *Plugin) TranscribeWithRetry(audioPath string, channelID string) (string, error) {
      maxRetries := 3
      
      for attempt := 1; attempt <= maxRetries; attempt++ {
          transcript, err := p.TranscribeAudio(audioPath)
          
          if err == nil {
              return transcript, nil
          }
          
          // Post retry notification
          if attempt < maxRetries {
              p.postMessage(channelID, fmt.Sprintf("⚠️ Transcription failed, retrying... (Attempt %d/%d)", attempt, maxRetries))
              time.Sleep(time.Duration(attempt*2) * time.Second) // Exponential backoff
          }
      }
      
      return "", fmt.Errorf("transcription failed after %d attempts", maxRetries)
  }
  ```

#### 6. Summarization

**Function**: `GenerateSummary()`
- **Location**: `summarization.go`
- **Purpose**: Extract action items and summary from transcript
- **Parameters**:
  - `transcript string`
- **Returns**: `(*Summary, error)`
- **Implementation**:
  ```go
  type Summary struct {
      ActionItems []string `json:"action_items"`
      Summary     string   `json:"summary"`
      Topics      []string `json:"topics"`
  }
  
  func (p *Plugin) GenerateSummary(transcript string) (*Summary, error) {
      prompt := fmt.Sprintf(`Analyze this meeting transcript and extract:

1. Action Items: Specific tasks mentioned (with assignees if mentioned)
2. Summary: 2-3 sentence overview of what was discussed
3. Topics: Key topics discussed (max 5)

Return ONLY valid JSON in this exact format:
{
  "action_items": ["@person to do X by Y", "Do Z (unassigned)"],
  "summary": "Brief summary here...",
  "topics": ["Topic 1", "Topic 2"]
}

Transcript:
%s`, transcript)
      
      client := resty.New()
      resp, err := client.R().
          SetHeader("Authorization", "Bearer "+p.getOpenAIKey()).
          SetHeader("Content-Type", "application/json").
          SetBody(map[string]interface{}{
              "model": "gpt-4o",
              "messages": []map[string]string{
                  {"role": "system", "content": "You are a meeting assistant that extracts action items and summaries."},
                  {"role": "user", "content": prompt},
              },
              "temperature": 0.3,
          }).
          Post("https://api.openai.com/v1/chat/completions")
      
      if err != nil {
          return nil, err
      }
      
      // Parse response
      var result struct {
          Choices []struct {
              Message struct {
                  Content string `json:"content"`
              } `json:"message"`
          } `json:"choices"`
      }
      
      err = json.Unmarshal(resp.Body(), &result)
      if err != nil {
          return nil, err
      }
      
      // Extract JSON from response
      content := result.Choices[0].Message.Content
      
      // Strip markdown code blocks if present
      content = strings.TrimPrefix(content, "```json\n")
      content = strings.TrimSuffix(content, "\n```")
      
      var summary Summary
      err = json.Unmarshal([]byte(content), &summary)
      if err != nil {
          return nil, err
      }
      
      return &summary, nil
  }
  ```

#### 7. Result Formatting & Posting

**Function**: `FormatAndPostResult()`
- **Location**: `formatter.go`
- **Purpose**: Format summary as markdown and post
- **Implementation**:
  ```go
  func (p *Plugin) FormatAndPostResult(recording *Recording, transcript string, summary *Summary) error {
      // Calculate duration
      duration := time.Since(recording.StartTime)
      
      // Get participants
      participants := p.getCallParticipants(recording.CallID)
      participantMentions := ""
      for _, user := range participants {
          participantMentions += "@" + user.Username + ", "
      }
      participantMentions = strings.TrimSuffix(participantMentions, ", ")
      
      // Format action items
      actionItemsText := ""
      for _, item := range summary.ActionItems {
          actionItemsText += fmt.Sprintf("- %s\n", item)
      }
      
      // Build message
      message := fmt.Sprintf(`## 🎯 Meeting Summary

**Duration**: %s  
**Participants**: %s  
**Started**: %s  
**Ended**: %s

---

### ✅ Action Items

%s

---

### 📝 Summary

%s

---

### 📄 Full Transcript

<details>
<summary>Click to expand full transcript</summary>

%s

</details>`, 
          formatDuration(duration),
          participantMentions,
          recording.StartTime.Format("3:04 PM"),
          time.Now().Format("3:04 PM"),
          actionItemsText,
          summary.Summary,
          transcript,
      )
      
      // Determine where to post
      channelID := recording.ChannelID
      if recording.IsDM {
          // DM all participants
          for _, user := range participants {
              p.postDM(user.Id, message)
          }
      } else {
          // Post in channel
          p.postMessage(channelID, message)
      }
      
      // Store in database
      p.storeTranscript(recording.CallID, transcript, summary)
      
      return nil
  }
  ```

#### 8. Database Operations

**Function**: `storeTranscript()`
- **Location**: `storage.go`
- **Purpose**: Save transcript to database
- **Implementation**:
  ```go
  func (p *Plugin) storeTranscript(callID string, transcript string, summary *Summary) error {
      // Get user's retention setting
      retentionDays := p.getUserRetentionSetting() // Default: 30 days
      expiresAt := time.Now().AddDate(0, 0, retentionDays).Unix()
      
      // If user set "never delete", expiresAt = 0
      if retentionDays == -1 {
          expiresAt = 0
      }
      
      // Insert transcript
      _, err := p.API.Store.Exec(`
          INSERT INTO meeting_transcripts 
          (id, call_id, transcript, created_at, expires_at)
          VALUES (?, ?, ?, ?, ?)`,
          model.NewId(),
          callID,
          transcript,
          time.Now().Unix(),
          expiresAt,
      )
      
      if err != nil {
          return err
      }
      
      // Insert summary
      summaryJSON, _ := json.Marshal(summary)
      _, err = p.API.Store.Exec(`
          INSERT INTO meeting_summaries
          (id, call_id, summary_data, created_at)
          VALUES (?, ?, ?, ?)`,
          model.NewId(),
          callID,
          string(summaryJSON),
          time.Now().Unix(),
      )
      
      return err
  }
  ```

#### 9. Cleanup Job

**Function**: `CleanupExpiredTranscripts()`
- **Location**: `storage.go`
- **Purpose**: Delete transcripts past their retention period
- **Implementation**:
  ```go
  func (p *Plugin) CleanupExpiredTranscripts() {
      // Run daily
      ticker := time.NewTicker(24 * time.Hour)
      
      for range ticker.C {
          now := time.Now().Unix()
          
          // Delete expired transcripts
          _, err := p.API.Store.Exec(`
              DELETE FROM meeting_transcripts
              WHERE expires_at > 0 AND expires_at < ?`,
              now,
          )
          
          if err != nil {
              p.API.LogError("Failed to cleanup transcripts", "error", err.Error())
          }
      }
  }
  ```

---

## 5. Database Schema

### Tables

#### Table: `meeting_recordings`
**Purpose**: Track active recordings

```sql
CREATE TABLE meeting_recordings (
    id VARCHAR(26) PRIMARY KEY,
    call_id VARCHAR(26) NOT NULL,
    channel_id VARCHAR(26) NOT NULL,
    started_by VARCHAR(26) NOT NULL,       -- user_id who started recording
    started_at BIGINT NOT NULL,
    ended_at BIGINT,
    status VARCHAR(20) NOT NULL,           -- 'recording', 'processing', 'completed', 'failed'
    audio_file_path TEXT,                  -- temp path before processing
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    
    INDEX idx_call_id (call_id),
    INDEX idx_status (status)
);
```

#### Table: `meeting_transcripts`
**Purpose**: Store transcripts with retention policy

```sql
CREATE TABLE meeting_transcripts (
    id VARCHAR(26) PRIMARY KEY,
    call_id VARCHAR(26) NOT NULL,
    recording_id VARCHAR(26) NOT NULL,     -- FK to meeting_recordings
    transcript TEXT NOT NULL,              -- Full transcript
    word_count INTEGER,
    created_at BIGINT NOT NULL,
    expires_at BIGINT,                     -- 0 = never expire, else unix timestamp
    
    INDEX idx_call_id (call_id),
    INDEX idx_expires_at (expires_at),
    FOREIGN KEY (recording_id) REFERENCES meeting_recordings(id)
);
```

#### Table: `meeting_summaries`
**Purpose**: Store extracted summaries and action items

```sql
CREATE TABLE meeting_summaries (
    id VARCHAR(26) PRIMARY KEY,
    call_id VARCHAR(26) NOT NULL,
    recording_id VARCHAR(26) NOT NULL,
    transcript_id VARCHAR(26) NOT NULL,
    summary TEXT NOT NULL,                 -- 2-3 sentence summary
    action_items JSONB,                    -- Array: ["item 1", "item 2"]
    topics JSONB,                          -- Array: ["topic 1", "topic 2"]
    created_at BIGINT NOT NULL,
    
    INDEX idx_call_id (call_id),
    FOREIGN KEY (recording_id) REFERENCES meeting_recordings(id),
    FOREIGN KEY (transcript_id) REFERENCES meeting_transcripts(id)
);
```

#### Table: `plugin_config`
**Purpose**: Store plugin configuration (API keys, settings)

```sql
CREATE TABLE plugin_config (
    key VARCHAR(255) PRIMARY KEY,
    value TEXT NOT NULL,
    encrypted BOOLEAN DEFAULT FALSE,       -- For API keys
    updated_at BIGINT NOT NULL
);
```

#### Table: `user_settings`
**Purpose**: Per-user configuration (retention preferences)

```sql
CREATE TABLE user_settings (
    user_id VARCHAR(26) PRIMARY KEY,
    retention_days INTEGER DEFAULT 30,    -- -1 = never delete, 0 = delete immediately
    auto_join BOOLEAN DEFAULT FALSE,      -- Future: auto-join all calls
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);
```

### Indexes

```sql
-- For cleanup job
CREATE INDEX idx_transcripts_expiry ON meeting_transcripts(expires_at) WHERE expires_at > 0;

-- For searching transcripts
CREATE INDEX idx_transcripts_created ON meeting_transcripts(created_at DESC);

-- For finding user's recordings
CREATE INDEX idx_recordings_user ON meeting_recordings(started_by, created_at DESC);
```

---

## 6. API Integrations

### OpenAI Whisper API

**Endpoint**: `POST https://api.openai.com/v1/audio/transcriptions`

**Request**:
```http
POST /v1/audio/transcriptions HTTP/1.1
Host: api.openai.com
Authorization: Bearer sk-...
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary

------WebKitFormBoundary
Content-Disposition: form-data; name="file"; filename="meeting.wav"
Content-Type: audio/wav

[Binary audio data]
------WebKitFormBoundary
Content-Disposition: form-data; name="model"

whisper-1
------WebKitFormBoundary
Content-Disposition: form-data; name="response_format"

json
------WebKitFormBoundary--
```

**Response**:
```json
{
  "text": "This is the full transcript of the meeting..."
}
```

**Error Handling**:
- 400: Invalid audio format → Post user error
- 401: Invalid API key → Log error, post admin message
- 429: Rate limit → Retry with backoff
- 500: Server error → Retry up to 3 times

---

### OpenAI GPT-4o API

**Endpoint**: `POST https://api.openai.com/v1/chat/completions`

**Request**:
```json
{
  "model": "gpt-4o",
  "messages": [
    {
      "role": "system",
      "content": "You are a meeting assistant that extracts action items and summaries."
    },
    {
      "role": "user",
      "content": "Analyze this meeting transcript...\n\n[transcript]"
    }
  ],
  "temperature": 0.3,
  "max_tokens": 2000
}
```

**Response**:
```json
{
  "choices": [
    {
      "message": {
        "content": "{\"action_items\": [...], \"summary\": \"...\", \"topics\": [...]}"
      }
    }
  ],
  "usage": {
    "prompt_tokens": 5000,
    "completion_tokens": 200,
    "total_tokens": 5200
  }
}
```

**Cost Calculation**:
- Input: $0.0025 per 1K tokens
- Output: $0.0025 per 1K tokens
- Average 30-min meeting: ~8K tokens input + 300 tokens output = $0.02 per meeting

---

### Mattermost Plugin API (Internal)

**Key API Methods**:

```go
// Post a message
p.API.CreatePost(&model.Post{
    UserId: p.botUserID,
    ChannelId: channelID,
    Message: "🤖 Recording started",
})

// Get channel info
channel, err := p.API.GetChannel(channelID)

// Get user info
user, err := p.API.GetUser(userID)

// Send DM
dmChannel, err := p.API.GetDirectChannel(userID1, userID2)
p.API.CreatePost(&model.Post{
    UserId: p.botUserID,
    ChannelId: dmChannel.Id,
    Message: "Summary here...",
})

// Database operations
result, err := p.API.Store.Exec("INSERT INTO ...")
rows, err := p.API.Store.Query("SELECT ...")
```

---

## 7. Security & Permissions

### Authentication

**Bot User**:
- Created automatically on plugin activation
- System-generated credentials
- No login access (API-only user)
- Displayed as "Meeting Bot" in UI

**API Keys**:
- Stored in `plugin_config` table
- Encrypted at rest using Mattermost's encryption
- Only System Admins can view/modify
- Configured via System Console

### Authorization

**Who Can Use the Bot**:
- Any user who is in an active call
- No special permissions required (channel member = can record)
- Configurable: Admin can restrict to Team Admins only (post-MVP)

**Who Can Configure Settings**:
- System Admins: API keys, global settings
- Regular users: Personal settings (retention preference)

### Data Privacy

**Audio Files**:
- Saved to `/tmp/` (ephemeral storage)
- Deleted immediately after transcription
- Never stored in database
- Not accessible via API

**Transcripts**:
- Stored in database with retention policy
- User-configurable: 30 days (default), forever, or immediate delete
- No encryption (stored as plain text for search)
- Accessible only to call participants (post-MVP: add access control)

**API Keys**:
- Never logged
- Never sent to client
- Encrypted in database
- Rotatable via System Console

### Compliance

**GDPR Considerations**:
- Users can delete their transcripts on demand
- Retention policies respect data minimization
- Audit log of who accessed what (post-MVP)

**Audio Recording Laws**:
- Bot announces when recording starts (informed consent)
- Recording is opt-in (user triggers it)
- No silent/automatic recording

---

## 8. Testing Strategy

### Unit Tests (Go)

**Files to Test**:

1. **transcription_test.go**:
```go
func TestTranscribeAudio(t *testing.T) {
    // Mock Whisper API response
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
        json.NewEncoder(w).Encode(map[string]string{
            "text": "Test transcript",
        })
    }))
    defer server.Close()
    
    plugin := &Plugin{/* ... */}
    transcript, err := plugin.TranscribeAudio("test.wav")
    
    assert.NoError(t, err)
    assert.Equal(t, "Test transcript", transcript)
}

func TestTranscribeWithRetry(t *testing.T) {
    // Test retry logic with failing API
    // ...
}
```

2. **summarization_test.go**:
```go
func TestGenerateSummary(t *testing.T) {
    // Mock GPT-4o response
    // ...
}
```

3. **formatter_test.go**:
```go
func TestFormatResult(t *testing.T) {
    summary := &Summary{
        ActionItems: []string{"Item 1", "Item 2"},
        Summary: "Test summary",
    }
    
    result := formatMarkdown(summary, "transcript...")
    
    assert.Contains(t, result, "## 🎯 Meeting Summary")
    assert.Contains(t, result, "Item 1")
}
```

4. **storage_test.go**:
```go
func TestStoreTranscript(t *testing.T) {
    // Test database insertion
    // ...
}

func TestCleanupExpiredTranscripts(t *testing.T) {
    // Test deletion of expired records
    // ...
}
```

### Integration Tests

**Test Scenarios**:

1. **End-to-End Recording Flow**:
   - Start recording via slash command
   - Simulate call end event
   - Verify transcript posted to channel
   - Verify data in database

2. **API Failure Handling**:
   - Mock Whisper API failure
   - Verify retry logic works
   - Verify error message posted

3. **Multiple Concurrent Recordings**:
   - Start 2 recordings in different channels
   - Verify both complete successfully
   - No data mixing between recordings

### Manual Testing Checklist

**Pre-Test Setup**:
- [ ] Configure OpenAI API key in System Console
- [ ] Install Calls plugin
- [ ] Create test audio file (5 minutes, clear speech)

**Test Cases**:

1. **Happy Path**:
   - [ ] Start voice call in test channel
   - [ ] Run `/meeting-bot start`
   - [ ] Verify bot joins (confirmation message)
   - [ ] Talk for 2-3 minutes
   - [ ] End call
   - [ ] Verify "Processing..." message appears
   - [ ] Verify formatted summary appears within 2 minutes
   - [ ] Verify action items, summary, and transcript sections
   - [ ] Verify transcript is collapsible

2. **No Active Call**:
   - [ ] Run `/meeting-bot start` with no call active
   - [ ] Verify error message: "No active call detected"

3. **Already Recording**:
   - [ ] Start recording in call
   - [ ] Try to start again
   - [ ] Verify message: "Already recording this call"

4. **API Failure**:
   - [ ] Temporarily disable API key
   - [ ] Record a call
   - [ ] Verify retry messages appear
   - [ ] Verify final error message after 3 retries

5. **Group DM Call**:
   - [ ] Start call in group DM
   - [ ] Record meeting
   - [ ] Verify summary posted in group DM

6. **User Settings**:
   - [ ] Run `/meeting-bot settings`
   - [ ] Change retention to "never delete"
   - [ ] Verify setting saved
   - [ ] Check database: expires_at = 0

7. **Long Meeting**:
   - [ ] Record 30+ minute call
   - [ ] Verify transcription completes
   - [ ] Verify no truncation

---

## 9. Error Handling

### Error Categories

**1. User Errors** (inform user, don't retry):
- No active call when starting recording
- Bot already recording this call
- Invalid command syntax

**2. Transient Errors** (retry):
- API timeout (Whisper/GPT-4o)
- Network connection lost
- Rate limit exceeded (429)

**3. Configuration Errors** (alert admin):
- API key not configured
- API key invalid/expired
- Database connection failed

**4. System Errors** (log, alert admin):
- WebRTC connection failed
- Audio encoding failed
- Disk full (can't save recording)

### Retry Strategy

**Exponential Backoff**:
```
Attempt 1: Wait 2 seconds
Attempt 2: Wait 4 seconds
Attempt 3: Wait 8 seconds
Max retries: 3
```

**Retry Logic**:
```go
func retryWithBackoff(fn func() error, maxAttempts int) error {
    for attempt := 1; attempt <= maxAttempts; attempt++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        if attempt < maxAttempts {
            backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
            time.Sleep(backoff)
        }
    }
    return fmt.Errorf("failed after %d attempts", maxAttempts)
}
```

### Error Messages (User-Facing)

```markdown
✅ Success: "🤖 Meeting Assistant joined the call and is recording"
✅ Success: "🤖 Processing your meeting... This may take 2 minutes."
✅ Success: [Formatted summary posted]

⚠️ Warning: "⚠️ Transcription attempt 1 failed, retrying... (Attempt 1/3)"
⚠️ Warning: "⚠️ No speech detected in recording. Please check audio setup."

❌ Error: "❌ No active call detected in this channel. Please join a call first."
❌ Error: "❌ Transcription failed after 3 attempts. Please contact your admin."
❌ Error: "❌ API key not configured. Admin: Configure in System Console > Plugins > Meeting Bot"
```

---

## 10. Risks & Unknowns

### Technical Risks

**Risk 1: Calls Plugin Integration Complexity**
- **Risk**: Calls plugin may not expose APIs we need for bot to join
- **Likelihood**: Medium
- **Impact**: High (blocks entire feature)
- **Mitigation**: 
  - Day 2: Investigate Calls plugin architecture (2-3 hours)
  - If Option A (extend plugin) is blocked, pivot to Option B (standalone WebRTC)
  - Option B is more complex but fully in our control

**Risk 2: WebRTC Audio Quality**
- **Risk**: Audio capture might be poor quality (echo, noise, low bitrate)
- **Likelihood**: Medium
- **Impact**: Medium (bad transcripts)
- **Mitigation**:
  - Test early with real calls (Day 3)
  - Implement audio preprocessing (noise reduction, normalization)
  - Document audio quality requirements for users

**Risk 3: API Cost Overruns**
- **Risk**: If many users record long meetings, API costs could be high
- **Likelihood**: Low (you have free API access)
- **Impact**: Low (but could be issue post-MVP)
- **Mitigation**:
  - Add rate limiting (10 recordings/user/day)
  - Add meeting length limit (60 minutes max)
  - Track usage metrics

**Risk 4: Whisper API Rate Limits**
- **Risk**: 50 requests/minute limit could be hit with concurrent recordings
- **Likelihood**: Low (small team)
- **Impact**: Medium (recordings fail)
- **Mitigation**:
  - Implement queue system for transcription
  - Retry with exponential backoff
  - Consider local Whisper as fallback

### Unknowns (To Investigate Day 2)

**Unknown 1: Calls Plugin Architecture**
- **Question**: Can we programmatically join calls as a bot user?
- **Investigation**: Read Calls plugin source code, check API
- **Time**: 2-3 hours
- **Backup plan**: Standalone WebRTC implementation (Option B)

**Unknown 2: Audio Format**
- **Question**: What format does Calls plugin use for audio? (Opus? PCM?)
- **Investigation**: Test with live call, inspect RTP packets
- **Time**: 1 hour
- **Backup plan**: Support multiple formats, convert as needed

**Unknown 3: Multi-Participant Audio Mixing**
- **Question**: Do we receive pre-mixed audio, or separate streams per participant?
- **Investigation**: Join test call, examine audio tracks
- **Time**: 1 hour
- **Impact on design**: Affects recording implementation

**Unknown 4: Call End Detection**
- **Question**: Does Calls plugin fire reliable `call_ended` event?
- **Investigation**: Review Calls plugin code, test with calls
- **Time**: 30 minutes
- **Backup plan**: Timeout-based detection (no audio for 30 seconds)

### Decision Points (Day 2)

After investigating Calls plugin:

**Decision 1: Integration Approach**
- [ ] Option A: Extend Calls plugin ✅ (if APIs exist)
- [ ] Option B: Standalone WebRTC (if Option A blocked)

**Decision 2: Audio Handling**
- [ ] Receive mixed audio (simpler)
- [ ] Receive separate streams (need to mix ourselves)

**Decision 3: Speaker Labels**
- [ ] Skip for MVP (confirmed) ✅
- [ ] Add post-MVP using Deepgram

---

## 11. Rollout Plan

### Phase 1: MVP (Days 2-6)

**Day 2: Investigation & Setup**
- ✅ Investigate Calls plugin architecture (2-3 hours)
- ✅ Set up plugin skeleton (1 hour)
- ✅ Implement slash command handler (1 hour)
- ✅ Create bot user (1 hour)
- ✅ Test command registration (30 min)

**Day 3: WebRTC & Recording**
- ✅ Implement WebRTC join logic (4 hours)
- ✅ Test joining live call (2 hours)
- ✅ Implement audio capture and buffering (2 hours)

**Day 4: Transcription & Summarization**
- ✅ Implement Whisper API integration (2 hours)
- ✅ Implement retry logic (1 hour)
- ✅ Implement GPT-4o summarization (2 hours)
- ✅ Test end-to-end flow with test audio (2 hours)

**Day 5: Formatting & Database**
- ✅ Implement markdown formatting (2 hours)
- ✅ Implement result posting logic (1 hour)
- ✅ Create database schema (1 hour)
- ✅ Implement storage functions (2 hours)
- ✅ Test with multiple recordings (2 hours)

**Day 6: Error Handling & Testing**
- ✅ Implement error handling and retries (2 hours)
- ✅ Add user settings (retention config) (2 hours)
- ✅ Manual testing with real calls (3 hours)
- ✅ Bug fixes (2 hours)

**Day 7: Polish & Demo Prep**
- ✅ Code cleanup and comments (1 hour)
- ✅ Write documentation (README) (2 hours)
- ✅ Create demo recording (1 hour)
- ✅ Practice demo (1 hour)
- ✅ Buffer for unexpected issues (3 hours)

### Phase 2: Enhancements (Post-MVP)

**If Time Permits (Days 7+)**:
- Add "Assistant!" button in Calls UI
- Implement speaker labels (Deepgram)
- Add timestamps to transcript
- Implement transcript editing
- Add meeting analytics

**Future Iterations**:
- Real-time transcription display (live captions)
- Meeting search across all transcripts
- Integration with calendar (auto-record scheduled meetings)
- Custom action item assignment (drag-and-drop)
- Export to PDF/Google Docs

---

## 12. Open Questions

### Pre-Development (Resolve on Day 2)

- [ ] **Calls Plugin API**: What APIs does Calls plugin expose?
- [ ] **Audio Format**: What codec/format is audio in? Opus? PCM?
- [ ] **Multi-Participant Mixing**: Pre-mixed or separate streams?
- [ ] **Call End Event**: Reliable `call_ended` event exists?

### During Development (Resolve as needed)

- [ ] **Performance**: Can we handle 10+ concurrent recordings?
- [ ] **Audio Quality**: Is captured audio good enough for Whisper?
- [ ] **Transcript Length**: Do we need to chunk very long meetings?
- [ ] **Bot Presence**: Should bot appear in participant list?

### Post-MVP (Nice to Have)

- [ ] Should we support video recording?
- [ ] Should we integrate with Google Calendar?
- [ ] Should we add meeting analytics?
- [ ] Should we allow custom prompts for summarization?

---

## 13. References

### Mattermost Documentation
- **Plugin Developer Guide**: https://developers.mattermost.com/integrate/plugins/
- **Plugin API Reference**: https://developers.mattermost.com/integrate/reference/server/server-reference/
- **Calls Plugin Repo**: https://github.com/mattermost/mattermost-plugin-calls
- **Plugin Best Practices**: https://developers.mattermost.com/integrate/plugins/best-practices/

### External Documentation
- **Pion WebRTC**: https://github.com/pion/webrtc
- **Pion WebRTC Examples**: https://github.com/pion/webrtc/tree/master/examples
- **OpenAI Whisper API**: https://platform.openai.com/docs/guides/speech-to-text
- **OpenAI GPT-4o API**: https://platform.openai.com/docs/api-reference/chat
- **Go Audio Processing**: https://github.com/go-audio/audio

### Related Code Examples
- **Mattermost Plugin Starter**: https://github.com/mattermost/mattermost-plugin-starter-template
- **Calls Plugin Source**: https://github.com/mattermost/mattermost-plugin-calls
- **Bot User Examples**: Search Mattermost codebase for `CreateBot()`

---

## Change Log

### 2025-01-15 - Initial Draft
- Created comprehensive TDD for AI Meeting Assistant Bot
- Defined MVP scope: Slash command trigger, batch transcription, no speaker labels
- Selected tech stack: Whisper (transcription), GPT-4o (summarization)
- Outlined 7-day development plan
- Documented all API integrations, database schema, error handling
- Identified risks and unknowns (Calls plugin integration)

---

## Next Steps

1. ✅ Review this TDD with team/stakeholders
2. ⏭️ Generate context file for Calls plugin and WebRTC (Day 2 AM)
3. ⏭️ Investigate Calls plugin architecture (Day 2 AM)
4. ⏭️ Create tasks.md breakdown based on this TDD (Day 2)
5. ⏭️ Begin implementation (Day 2 PM)

---

**STATUS**: ✅ TDD Complete - Ready for Development