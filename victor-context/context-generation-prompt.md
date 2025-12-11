# Context File Generation

## Important: Iterative Approach

**We ARE creating all the sections** (the structure exists from the template).

**But we fill them in WHEN WE GET THERE / AS NEEDED** during implementation, not upfront.

As you work through tasks, fill in the relevant sections when you encounter them. Don't try to predict everything - let the implementation guide what goes in context.

---

I'm creating a semantic context file to help you understand the Mattermost 
codebase for building an AI Meeting Assistant Bot plugin.

I have a template at context/features/meeting-bot/context.md

Your task: Fill in this template following these rules (as you encounter each section):

1. FEATURE OVERVIEW section:
   - Explain we're building a bot that joins voice calls, transcribes audio, 
     and posts summaries
   - Keep it to 3-4 sentences

2. ARCHITECTURE MAP section:
   - Identify which Mattermost systems we'll interact with:
     * Plugin system (for bot logic)
     * WebSocket/WebRTC (for joining calls)
     * Database (for storing transcripts)
     * HTTP API (for posting messages)
   - Draw a simple data flow diagram in text

3. CODEBASE LOCATIONS section:
   - Search the Mattermost repo for:
     * Plugin-related files (server/plugin/)
     * WebSocket handling (server/app/web_*.go)
     * Post creation (server/app/post.go)
   - List top 10 most relevant files with 1-sentence purpose each
   - DO NOT paste entire file contents

4. CODE PATTERNS section:
   - Find examples of:
     * How plugins register hooks (show short code snippet)
     * How bots post messages (show API call pattern)
     * How WebSocket events work (show event structure)
   - Keep code snippets under 15 lines

5. KEY FUNCTIONS section:
   - Identify these functions and document them:
     * CreatePost() - for bot to send messages
     * Plugin.OnActivate() - for initializing plugin
     * WebSocket broadcast functions
   - For each: location, purpose, parameters, returns

6. DATABASE SCHEMA section:
   - Identify tables we'll need:
     * Posts (for bot messages)
     * Users (bot user account)
     * Any new tables for transcripts
   - Document relevant columns only

7. EXTERNAL DOCUMENTATION section:
   - Add links to:
     * Mattermost plugin developer docs
     * Pion WebRTC library (we'll use this)
     * OpenAI Whisper API docs
     * Any other relevant external resources

Use the existing @codebase knowledge. Search intelligently. 
Prioritize semantic understanding over completeness.

Fill in sections as you work through tasks - don't try to do everything upfront.
```

#### Step 3: Iterative refinement (2-3 hours)

Work through each section with Cursor:
- **Review each section** before moving to next
- **Ask Cursor to clarify** if something is unclear
- **Add your own notes** as you learn
- **Keep snippets small** (under 15 lines)

#### Step 4: Test the context file (30 min)

Open a **new Cursor tab** (fresh context):
- Load the context file
- Ask Cursor: "Where should I start building the meeting bot plugin?"
- See if Cursor gives coherent, specific answer
- If yes → strategy works! If no → refine context file

---

### Phase 2: Hierarchical Context (Add Core Layer)

After first context file works, create reusable core context.

**Files to create**:
```
context/
├── core/
│   ├── plugin-system.md       # How plugins work in general
│   ├── api-patterns.md         # REST API conventions
│   ├── database-access.md      # How to query database
│   └── websocket-system.md     # Real-time communication patterns
└── features/
    └── meeting-bot/
        └── context.md          # Feature-specific context
^ prompt

# Context File Usage Rules

## When Starting Work on a Feature

1. Load the appropriate context file FIRST
   - Core context: Always load all files in context/core/
   - Feature context: Load context/features/{feature-name}/context.md

2. Refer to context file for:
   - Architecture understanding
   - Code patterns and conventions
   - Relevant file locations
   - Common gotchas

3. Update context file when you learn something important:
   - Add to "WORKING NOTES" section immediately
   - Promote to main sections during daily review

## Context File Priority

Priority order when generating responses:
1. Context files (highest priority)
2. Actual codebase files
3. External documentation
4. General knowledge

## Keep Context Fresh

- Review context file at start of each day
- Update "LAST UPDATED" when making changes
- If context contradicts actual code, UPDATE CONTEXT (code is source of truth)