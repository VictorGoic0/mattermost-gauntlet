# Calls Plugin Investigation

**Goal**: Understand Calls plugin architecture and decide integration approach (Option A vs Option B)

**Date Started**: 2025-01-15

---

## Investigation Tasks

- [x] 1.1 Clone Calls plugin repository
- [x] 1.2 Read Calls plugin documentation
- [x] 1.3 Investigate WebRTC implementation
- [x] 1.4 Search for bot/recording APIs
- [x] 1.5 Examine audio handling
- [x] 1.6 Check for call lifecycle events
- [x] 1.7 Document findings (this file)
- [x] 1.8 If Option A: Document integration points (N/A - Option A not viable)
- [x] 1.9 If Option B: Research Pion WebRTC

**Decision**: ✅ **Option B (Standalone WebRTC)** - See Decision section below

---

## Findings

### Calls Plugin Repository
- **URL**: https://github.com/mattermost/mattermost-plugin-calls
- **Location**: `research/calls-plugin/`
- **Version**: Latest main branch (commit: c0cb2c605e)

### Architecture Overview

**Key Discovery**: Calls plugin uses an **external `rtcd` service** for WebRTC handling
- `rtcd` = RTC daemon (separate service)
- Plugin communicates with rtcd via client interface (`interfaces.RTCDClient`)
- WebRTC connections are managed by rtcd, not directly by the plugin

**Structure**:
- `server/rtcd.go` - Manages connection to rtcd service
- `server/public/` - Public types and interfaces
- `server/api.go` - HTTP API endpoints
- `server/bot_api.go` - Bot-specific APIs (for Calls plugin's own bot)
- `server/recording_api.go` - Recording job management

### WebRTC Implementation

**Finding**: WebRTC is handled by external `rtcd` service, not directly in plugin
- Plugin uses `github.com/mattermost/rtcd/service` package
- Plugin creates RTCD client connections to rtcd hosts
- Actual WebRTC signaling/connection happens in rtcd service

**Key Files**:
- `server/rtcd.go` - RTCD client manager, connection handling
- Uses `interfaces.RTCDClient` interface to communicate with rtcd

### Available APIs

**HTTP API Endpoints** (from `api_router.go`):
- `/calls/{call_id}/active` - Get active call
- `/calls/{call_id}/recording/{action}` - Start/stop recording (requires host permissions)
- `/bot/*` - Bot APIs (requires bot session authentication)
  - `/bot/calls/{call_id}/recordings` - Post recording info
  - `/bot/calls/{call_id}/transcriptions` - Post transcription info
  - `/bot/calls/{call_id}/jobs/{job_id}/status` - Update job status

**Important**: Bot APIs are for Calls plugin's own bot, not external bots
- Bot session check: `isBotSession(r)` checks if request is from Calls plugin bot
- External plugins cannot use these endpoints

**No Public API for Programmatic Join**:
- No endpoint found for external plugins to join calls
- No API to get WebRTC signaling info for external use
- Recording uses external "job service" (`getJobService().RunJob()`)

### Audio Handling

**Finding**: Audio handling is done by rtcd service, not directly in Calls plugin
- Calls plugin sends/receives WebRTC signaling via WebSocket (`wsEventSignal`)
- Actual audio codec/format is handled by rtcd service
- Plugin receives SDP and ICE candidate messages via WebSocket
- Audio format is likely Opus (standard for WebRTC voice)

**WebSocket Signaling** (from `websocket.go`):
- `wsEventSignal` - WebRTC signaling messages (SDP, ICE candidates)
- `clientMessageTypeSDP` - SDP offer/answer messages
- `clientMessageTypeICE` - ICE candidate messages
- Signaling is relayed through Mattermost WebSocket, then to rtcd service

**Key Insight**: To join a call, we'd need to:
1. Get call info (call_id, channel_id)
2. Connect to rtcd service (or use WebRTC directly)
3. Exchange SDP offers/answers
4. Exchange ICE candidates
5. Receive audio tracks

**Challenge**: We don't have direct access to rtcd service or its signaling protocol

### Call Lifecycle Events

**WebSocket Events** (from `websocket.go`):
- `user_joined` - User joined call
- `user_left` - User left call
- `call_job_state` - Recording/transcription job state changes

**Plugin Communication**:
- Mattermost plugins CAN communicate via `PluginHTTP()` API
- Format: `/plugins/{plugin_id}/*` 
- However, Calls plugin's bot APIs require bot session authentication
- We could potentially call public endpoints like `/calls/{channel_id}/active`

**Question**: Are WebSocket events accessible to other plugins?
- Mattermost plugin system doesn't expose other plugins' WebSocket events directly
- We'd need to use HTTP APIs or database queries to detect call state changes

---

## Decision: Option A vs Option B

**Status**: ✅ **DECISION: Option B (Standalone WebRTC)**

### Option A: Extend Calls Plugin
- **Feasibility**: ❌ **NOT VIABLE**
- **Reason**: 
  - Calls plugin doesn't expose APIs for external plugins to join calls
  - Bot APIs are only for Calls plugin's own bot (requires bot session authentication)
  - WebRTC is handled by external `rtcd` service, not directly accessible
  - No public API to get WebRTC signaling info (SDP offers/answers)
  - Would require modifying Calls plugin codebase (not ideal for separate plugin)

### Option B: Standalone WebRTC ✅
- **Feasibility**: ✅ **VIABLE**
- **Approach**:
  1. Query for active calls via Mattermost database or HTTP API
  2. Use Pion WebRTC to create peer connection independently
  3. Need to figure out signaling (how to get SDP offer/answer)
  4. Join as bot user via WebRTC directly
- **Challenges**:
  - **Signaling**: How do we get WebRTC connection details? (SDP, ICE candidates)
    - Calls plugin uses rtcd service for signaling
    - We'd need to either: (a) connect to rtcd service, or (b) implement our own WebRTC signaling
  - **RTCD Service**: Calls plugin uses external rtcd service
    - rtcd is a separate service that handles WebRTC connections
    - We'd need to understand rtcd's API or implement standalone WebRTC
  - **Audio Format**: Likely Opus (standard WebRTC codec)
  - **Call Detection**: Can query via HTTP API (`/calls/{channel_id}/active`) or database
- **Next Steps**:
  - Research Pion WebRTC examples for joining existing sessions
  - Investigate rtcd service API (if accessible)
  - Consider alternative: Query database for active calls, implement standalone WebRTC join
  - Test if we can get TURN/STUN server config from Calls plugin

---

## Notes & Gotchas

### Key Discoveries
1. **Calls plugin uses external rtcd service** - WebRTC is not handled directly in plugin
2. **No public API for programmatic join** - Bot APIs are only for Calls plugin's own bot
3. **Plugin-to-plugin communication possible** - Can use `PluginHTTP()` to call Calls plugin endpoints
4. **WebSocket signaling** - SDP/ICE messages go through Mattermost WebSocket, then to rtcd
5. **Database access** - Can query Mattermost database directly for call info (if we have access)

### Challenges for Option B
1. **Signaling complexity** - Need to understand how to exchange SDP/ICE with rtcd or implement our own
2. **rtcd service** - May need to connect to rtcd service directly, or bypass it entirely
3. **Audio format** - Need to confirm Opus codec and sample rate
4. **TURN/STUN servers** - Need to get ICE server configuration

### Potential Approaches
1. **Query database for active calls** - Use Mattermost plugin API to query `calls` table
2. **Get TURN credentials** - Calls plugin has `/turn-credentials` endpoint (check if accessible)
3. **Standalone WebRTC** - Use Pion WebRTC to create peer connection independently
4. **Monitor call state** - Poll or watch for call start/end events via database or HTTP API

---

## Pion WebRTC Research (Task 1.9)

**Status**: ✅ Completed (initial research)

**Resources**:
- Pion WebRTC: https://github.com/pion/webrtc
- Examples: https://github.com/pion/webrtc/tree/master/examples

**Key Concepts for Joining Existing Session**:
1. Create PeerConnection with ICE servers (TURN/STUN)
2. Create offer or answer SDP
3. Exchange SDP with other participants (signaling)
4. Exchange ICE candidates
5. Receive audio/video tracks via `OnTrack` handler

**Challenge Identified**: Signaling mechanism
- Calls plugin uses rtcd service for signaling
- We'll need to implement our own signaling approach
- Options: (a) Connect to rtcd service if API exists, (b) Implement standalone WebRTC with custom signaling

**Next Steps for Implementation**:
- Research rtcd service API (if accessible)
- Implement Pion WebRTC peer connection
- Figure out how to exchange SDP/ICE with existing call participants
- Test audio capture from WebRTC tracks

---

## LAST UPDATED
2025-01-15 - Investigation in progress, Option B decision made

