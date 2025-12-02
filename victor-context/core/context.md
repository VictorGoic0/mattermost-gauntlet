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
[High-level: What are we building? Why? What's the end goal?]

---

## ARCHITECTURE MAP

### System Components Involved
[List the major components/systems this feature touches]

Example:
- Plugin System (for extending Mattermost)
- WebSocket Server (for real-time communication)
- Database Layer (for persistence)
- React Frontend (for UI)

### Data Flow
[Describe how data flows through the system for this feature]

Example:
User clicks button → Frontend sends HTTP request → Go handler processes 
→ Saves to database → Broadcasts WebSocket event → All clients update UI

---

## CODEBASE LOCATIONS

### Core Files (Must Understand)
[List files with brief purpose, NOT full content]

Example:
- `server/plugin/api.go` - Plugin API interface, defines hooks
- `server/app/post.go` - Message posting logic
- `webapp/src/actions/posts.js` - Redux actions for posts

### Related Files (May Need to Touch)
[Secondary files that might be relevant]

### External Dependencies
[Third-party libraries, APIs, services]

Example:
- Pion WebRTC (https://github.com/pion/webrtc) - for joining calls
- OpenAI Whisper API (https://platform.openai.com/docs/guides/speech-to-text)

---

## CODE PATTERNS & CONVENTIONS

### Pattern: [Name of Pattern]
**Where it's used**: [Location]
**What it does**: [Explanation]
**Example**:
```go
[Small, relevant code snippet showing the pattern]
```

### Gotcha: [Name of Gotcha]
**What**: [What's the issue]
**Why**: [Why it happens]
**How to avoid**: [Solution]

---

## KEY FUNCTIONS & APIS

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
[Add notes as you discover things]

### Questions / Unknowns
[Things you're still figuring out]

### Debugging Tips
[Tricks that helped during development]

---

## LAST UPDATED
[Date and what changed]