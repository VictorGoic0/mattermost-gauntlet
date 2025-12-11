# tasks-3.md

## Overview
This file covers Result Formatting, Error Handling, User Settings, Cleanup, and Demo Preparation.

**PRs Covered**: #5-8
**Days**: Day 5-7
**Estimated Time**: ~24 hours total

---

## PR #5: Result Formatting & Posting

**Branch**: `feature/meeting-bot-formatting`  
**Description**: Format meeting summaries as markdown and post to appropriate channels/DMs.

**Dependencies**: PR #4 merged (summarization working)

**Success Criteria**:
- [ ] Summaries formatted as rich markdown
- [ ] Action items displayed prominently
- [ ] Transcript collapsible
- [ ] Correct delivery (channel vs DM)
- [ ] @mentions work in action items
- [ ] Duration and timestamps calculated correctly
- [ ] All formatting tests pass

---

### Markdown Formatting Implementation

**Goal**: Create beautiful formatted message from summary data.

**Reference**: 
- 📄 TDD Section 2 (User Experience - Result Format)
- 📄 TDD Section 4 (Implementation Details - formatter.go)

#### Tasks:

- [ ] 5.1 Create formatter.go file
  - Add to `server/plugins/meeting-bot/formatter.go`

- [ ] 5.2 Implement FormatResult() function signature
  ```go
  func (p *Plugin) FormatResult(recording *Recording, transcript string, summary *Summary) string
  ```

- [ ] 5.3 Calculate meeting duration
  - duration = recording.EndTime - recording.StartTime
  - Format as "X minutes" or "X hours Y minutes"

- [ ] 5.4 Format start time
  - Convert timestamp to "3:04 PM" format
  - Use user's timezone (future enhancement: default UTC)

- [ ] 5.5 Format end time
  - Same format as start time

- [ ] 5.6 Get participant list
  - Call getCallParticipants(recording.CallID)
  - Format as "@username, @username, @username"

- [ ] 5.7 Build header section
  ```markdown
  ## 🎯 Meeting Summary
  
  **Duration**: X minutes
  **Participants**: @user1, @user2, @user3
  **Started**: 2:15 PM
  **Ended**: 2:47 PM
  
  ---
  ```

- [ ] 5.8 Format action items section
  - Loop through summary.ActionItems
  - Prefix each with "- "
  - If empty, show "No action items identified"

- [ ] 5.9 Build action items section
  ```markdown
  ### ✅ Action Items
  
  - @alice to send proposal by Friday
  - @bob to review budget numbers
  - Research competitor pricing (unassigned)
  
  ---
  ```

- [ ] 5.10 Build summary section
  ```markdown
  ### 📝 Summary
  
  {summary.Summary}
  
  ---
  ```

- [ ] 5.11 Build collapsible transcript section
  ```markdown
  ### 📄 Full Transcript
  
  <details>
  <summary>Click to expand full transcript</summary>
  
  {transcript}
  
  </details>
  ```

- [ ] 5.12 Concatenate all sections
  - Combine header + action items + summary + transcript
  - Return full markdown string

- [ ] 5.13 Handle empty transcript
  - If transcript is empty, show "Transcript unavailable"
  - Don't show empty <details> block

- [ ] 5.14 Handle very long transcripts
  - If transcript > 40,000 chars, truncate
  - Add note: "(Transcript truncated due to length)"

---

### Delivery Logic Implementation

**Goal**: Determine where to post results (channel vs DMs).

**Reference**: 
- 📄 TDD Section 2 (User Experience - Result Delivery Logic)

#### Tasks:

- [ ] 5.15 Implement determineDeliveryTarget() function
  ```go
  func (p *Plugin) determineDeliveryTarget(recording *Recording) DeliveryTarget
  ```

- [ ] 5.16 Check if recording was in channel
  - If recording.ChannelID is a public/private channel
  - Return: Post to that channel

- [ ] 5.17 Check if recording was in group DM
  - If recording.ChannelID is a group DM (3+ people)
  - Return: Post to that group DM

- [ ] 5.18 Check if recording was in 1-on-1 DM
  - If recording.ChannelID is a DM (2 people)
  - Return: Post to that DM

- [ ] 5.19 Handle edge cases
  - If channel was deleted, post to participants as DMs
  - If can't determine, default to posting in original channel

---

### Posting Implementation

**Goal**: Post formatted message to correct location.

**Reference**: 
- 📄 TDD Section 4 (Implementation Details - FormatAndPostResult)

#### Tasks:

- [ ] 5.20 Implement PostResult() function
  ```go
  func (p *Plugin) PostResult(recording *Recording, transcript string, summary *Summary) error
  ```

- [ ] 5.21 Format the result
  - Call FormatResult()
  - Get markdown string

- [ ] 5.22 Determine delivery target
  - Call determineDeliveryTarget()
  - Get target channel/DM

- [ ] 5.23 Post to channel (if target is channel)
  - Create Post object
  - Set ChannelId
  - Set UserId = bot user ID
  - Set Message = formatted markdown
  - Call p.API.CreatePost()

- [ ] 5.24 Post to group DM (if target is group DM)
  - Same as channel posting
  - ChannelId is the group DM channel ID

- [ ] 5.25 Post to individual DMs (if target is "DM all participants")
  - Loop through participants
  - For each user, get or create DM channel
  - Post message to each DM

- [ ] 5.26 Handle posting errors
  - If post fails, log error
  - Try alternative delivery (DM to admin)
  - Don't fail entire process

- [ ] 5.27 Update recording status in database
  - Set status = "completed"
  - Set updated_at timestamp

---

### Special Formatting Features

**Goal**: Add polish to formatted messages.

#### Tasks:

- [ ] 5.28 Parse @mentions in action items
  - Detect @username patterns
  - Convert to Mattermost mention format
  - Verify user exists

- [ ] 5.29 Add emojis to section headers
  - 🎯 for meeting summary
  - ✅ for action items
  - 📝 for summary
  - 📄 for transcript

- [ ] 5.30 Format duration intelligently
  - < 1 hour: "X minutes"
  - >= 1 hour: "X hours Y minutes"
  - Round to nearest minute

- [ ] 5.31 Highlight urgency in action items
  - Detect words: "urgent", "ASAP", "today", "tomorrow"
  - Add 🔥 emoji for urgent items

- [ ] 5.32 Add word count to transcript summary
  - Count words in transcript
  - Show: "Click to expand full transcript (X words)"

---

### Integration with Processing Flow

**Goal**: Call formatting after summarization completes.

**Reference**: 
- 📄 TDD Section 3 (Data Flow)

#### Tasks:

- [ ] 5.33 Update processRecording() function
  - After GenerateSummary() succeeds
  - Call PostResult()

- [ ] 5.34 Handle posting success
  - Log successful post
  - Mark recording complete in database

- [ ] 5.35 Handle posting failure
  - Log error with details
  - Retry once
  - If still fails, notify admin

---

### Testing

**Goal**: Verify formatting and posting work correctly.

#### Tasks:

- [ ] 5.36 Test formatting with real data
  - Use data from PRs #3-4
  - Generate formatted message
  - Verify markdown renders correctly

- [ ] 5.37 Test @mentions in action items
  - Include @username in test summary
  - Verify mention is clickable
  - Verify user gets notified

- [ ] 5.38 Test collapsible transcript
  - Click "expand full transcript"
  - Verify transcript shows
  - Click again to collapse

- [ ] 5.39 Test posting to channel
  - Record meeting in test channel
  - Verify result posted to same channel

- [ ] 5.40 Test posting to group DM
  - Record meeting in group DM
  - Verify result posted to group DM

- [ ] 5.41 Test posting to individual DMs
  - Simulate deleted channel
  - Verify DMs sent to all participants

- [ ] 5.42 Test with long transcript
  - Create 2-hour meeting transcript
  - Verify transcript truncates gracefully

- [ ] 5.43 Test with empty action items
  - Summary with no action items
  - Verify shows "No action items identified"

- [ ] 5.44 Verify formatting on mobile
  - Open Mattermost mobile app
  - View posted summary
  - Verify readable on small screen

---

### PR #5 Wrap-Up

- [ ] 5.45 Add formatting examples to README
  - Include screenshot or example markdown

- [ ] 5.46 Document delivery logic
  - Explain when posts go to channel vs DM

- [ ] 5.47 Commit and create PR
  - Git add all files
  - Commit: "feat: Implement result formatting and posting"
  - Push to branch
  - Create PR with formatted example

**PR #5 Complete** ✅

---

## PR #6: Error Handling & Retries

**Branch**: `feature/meeting-bot-error-handling`  
**Description**: Comprehensive error handling, retry logic, and user notifications.

**Dependencies**: PR #5 merged

**Success Criteria**:
- [ ] All error cases handled gracefully
- [ ] Users get clear error messages
- [ ] Retries work as expected
- [ ] No crashes from unexpected errors
- [ ] Admin alerted for config issues
- [ ] All error handling tests pass

---

### Comprehensive Error Handling

**Goal**: Handle all possible error scenarios.

**Reference**: 
- 📄 TDD Section 9 (Error Handling)

#### Tasks:

- [ ] 6.1 Create errors.go file
  - Add to `server/plugins/meeting-bot/errors.go`

- [ ] 6.2 Define custom error types
  ```go
  type MeetingBotError struct {
      Code    string
      Message string
      Err     error
  }
  ```

- [ ] 6.3 Define error codes
  - ErrNoActiveCall
  - ErrAlreadyRecording
  - ErrAPIKeyNotConfigured
  - ErrAPIKeyInvalid
  - ErrTranscriptionFailed
  - ErrSummarizationFailed
  - ErrPostingFailed
  - ErrDatabaseError

- [ ] 6.4 Implement error constructors
  - NewNoActiveCallError()
  - NewAPIKeyError()
  - Etc for each error type

- [ ] 6.5 Implement Error() method
  - Return formatted error message
  - Include error code and details

---

### User-Facing Error Messages

**Goal**: Clear, actionable error messages for users.

**Reference**: 
- 📄 TDD Section 9 (Error Handling - Error Messages)

#### Tasks:

- [ ] 6.6 Create error message mapping
  - Map error codes to user-friendly messages
  - Include emojis and formatting

- [ ] 6.7 Implement postUserError() helper
  ```go
  func (p *Plugin) postUserError(channelID string, err error)
  ```

- [ ] 6.8 Handle "no active call" error
  - Message: "❌ No active call detected in this channel. Please join a call first."
  - Post as ephemeral (only user sees)

- [ ] 6.9 Handle "already recording" error
  - Message: "🤖 Already recording this call (started by @username at HH:MM)"
  - Include who started and when

- [ ] 6.10 Handle "API key not configured" error
  - Message: "❌ OpenAI API key not configured. Admin: Configure in System Console > Plugins > Meeting Bot"
  - Post to channel (visible to all)
  - Also DM system admin

- [ ] 6.11 Handle "API key invalid" error
  - Message: "❌ OpenAI API key is invalid. Admin: Please check your API key."
  - Alert system admin via DM

- [ ] 6.12 Handle "transcription failed" error
  - Message: "❌ Transcription failed after 3 attempts. Audio has been saved for manual review."
  - Include recording ID for reference

- [ ] 6.13 Handle "posting failed" error
  - Message: "⚠️ Meeting processed successfully, but couldn't post result. Please check with admin."
  - Try to DM result to participants as fallback

---

### Retry Logic Enhancement

**Goal**: Robust retry logic with user feedback.

**Reference**: 
- 📄 TDD Section 9 (Error Handling - Retry Strategy)

#### Tasks:

- [ ] 6.14 Create retry.go file
  - Add to `server/plugins/meeting-bot/retry.go`

- [ ] 6.15 Implement generic retry function
  ```go
  func RetryWithBackoff(fn func() error, maxAttempts int, notifyFn func(attempt int)) error
  ```

- [ ] 6.16 Implement exponential backoff
  - Calculate delay: 2^attempt seconds
  - Max delay: 60 seconds

- [ ] 6.17 Add jitter to backoff
  - Add random ±20% to delay
  - Prevents thundering herd

- [ ] 6.18 Call notifyFn on each retry
  - Post retry message to channel
  - Include attempt number and reason

- [ ] 6.19 Update TranscribeWithRetry() to use new retry
  - Replace existing retry logic
  - Use RetryWithBackoff()

- [ ] 6.20 Update GPT-4o calls to use retry
  - Wrap API call in RetryWithBackoff()
  - Retry on rate limit or timeout

---

### Graceful Degradation

**Goal**: Partial success is better than total failure.

#### Tasks:

- [ ] 6.21 Allow transcript without summary
  - If summarization fails, still post transcript
  - Note: "Summary generation failed, showing transcript only"

- [ ] 6.22 Allow summary without topics
  - If GPT-4o doesn't return topics, continue
  - Show summary and action items

- [ ] 6.23 Allow empty action items
  - If no action items found, don't fail
  - Show message: "No action items identified"

- [ ] 6.24 Handle partial API responses
  - If API returns incomplete JSON, use what's available
  - Log warning about incomplete data

---

### Admin Notifications

**Goal**: Alert admins about configuration or system issues.

#### Tasks:

- [ ] 6.25 Implement getSystemAdmins() helper
  - Query Mattermost for system admin users
  - Return list of admin IDs

- [ ] 6.26 Implement notifyAdmins() function
  - Send DM to all system admins
  - Include error details and suggested fix

- [ ] 6.27 Notify admins on API key issues
  - When API key missing or invalid
  - Include link to System Console settings

- [ ] 6.28 Notify admins on repeated failures
  - If > 3 recordings fail in a row
  - Suggest checking API status or keys

- [ ] 6.29 Notify admins on disk space issues
  - If disk full when saving recordings
  - Suggest cleanup or expanding storage

---

### Error Recovery

**Goal**: Recover from errors when possible.

#### Tasks:

- [ ] 6.30 Implement manual retry command
  - `/meeting-bot retry {recording_id}`
  - Attempts to reprocess failed recording

- [ ] 6.31 Implement retry logic for failed recordings
  - Query database for failed recordings
  - Get audio file path (if still exists)
  - Retry transcription and summarization

- [ ] 6.32 Handle stale recordings
  - If recording > 24 hours old and status = "processing"
  - Mark as failed
  - Clean up resources

- [ ] 6.33 Implement cleanup for orphaned files
  - Check /tmp/ for old audio files
  - Delete files > 24 hours old
  - Log cleanup actions

---

### Testing

**Goal**: Test all error scenarios.

#### Tasks:

- [ ] 6.34 Test "no active call" error
  - Run command without call
  - Verify error message

- [ ] 6.35 Test "already recording" error
  - Start recording twice
  - Verify second attempt fails gracefully

- [ ] 6.36 Test API key errors
  - Remove API key
  - Try to record
  - Verify error messages and admin notification

- [ ] 6.37 Test transcription retry
  - Mock Whisper API to fail twice, succeed third time
  - Verify retry messages posted
  - Verify final success

- [ ] 6.38 Test graceful degradation
  - Mock GPT-4o to return partial response
  - Verify partial data still used

- [ ] 6.39 Test manual retry command
  - Create failed recording
  - Run retry command
  - Verify reprocessing works

---

### PR #6 Wrap-Up

- [ ] 6.40 Document all error codes
  - Create error reference in README

- [ ] 6.41 Add troubleshooting guide
  - Common errors and solutions

- [ ] 6.42 Commit and create PR
  - Git add all files
  - Commit: "feat: Comprehensive error handling and retry logic"
  - Push to branch
  - Create PR

**PR #6 Complete** ✅

---

## PR #7: User Settings & Cleanup

**Branch**: `feature/meeting-bot-settings`  
**Description**: User-configurable settings and automated cleanup jobs.

**Dependencies**: PR #6 merged

**Success Criteria**:
- [ ] Users can configure retention settings
- [ ] Settings stored and applied correctly
- [ ] Cleanup job deletes expired transcripts
- [ ] Settings accessible via command
- [ ] All settings tests pass

---

### User Settings Implementation

**Goal**: Allow users to configure personal preferences.

**Reference**: 
- 📄 TDD Section 5 (Database Schema - user_settings)

#### Tasks:

- [ ] 7.1 Create settings.go file
  - Add to `server/plugins/meeting-bot/settings.go`

- [ ] 7.2 Define UserSettings struct
  ```go
  type UserSettings struct {
      UserID        string
      RetentionDays int  // -1 = forever, 0 = immediate delete, N = N days
      CreatedAt     int64
      UpdatedAt     int64
  }
  ```

- [ ] 7.3 Implement getUserSettings() function
  - Query user_settings table
  - Return settings or default if not exists

- [ ] 7.4 Implement saveUserSettings() function
  - Insert or update user_settings
  - Validate retention value

- [ ] 7.5 Implement default settings
  - Default retention: 30 days
  - Applied when user has no saved settings

---

### Settings Command Implementation

**Goal**: `/meeting-bot settings` command to view/change settings.

**Reference**: 
- 📄 TDD Section 4 (Implementation Details - commands.go)

#### Tasks:

- [ ] 7.6 Update handleSettings() in commands.go
  - Remove stub implementation

- [ ] 7.7 Implement settings view
  - Show current settings
  - Format as ephemeral message

- [ ] 7.8 Show retention setting
  - Display current retention_days value
  - Format: "30 days", "Forever", "Immediate delete"

- [ ] 7.9 Show update instructions
  - Explain how to change settings
  - Usage: `/meeting-bot settings retention [days|forever|immediate]`

- [ ] 7.10 Implement settings update subcommand
  - Parse arguments
  - Validate new value

- [ ] 7.11 Handle retention setting update
  - `/meeting-bot settings retention 30` → 30 days
  - `/meeting-bot settings retention forever` → -1 (never delete)
  - `/meeting-bot settings retention immediate` → 0 (delete after processing)

- [ ] 7.12 Validate retention value
  - Must be: -1, 0, or positive integer
  - Max: 365 days (prevent unreasonable values)

- [ ] 7.13 Save updated settings
  - Call saveUserSettings()
  - Confirm to user: "✅ Settings updated"

- [ ] 7.14 Show warning for immediate delete
  - If user sets retention = 0
  - Warning: "⚠️ Transcripts will be deleted immediately after posting. This cannot be undone."

---

### Apply Retention Settings

**Goal**: Apply user's retention preference when storing transcripts.

**Reference**: 
- 📄 TDD Section 4 (Implementation Details - storeTranscript)

#### Tasks:

- [ ] 7.15 Update storeTranscript() in storage.go
  - Get user's retention setting
  - Calculate expires_at based on setting

- [ ] 7.16 Handle retention = 30 days (default)
  - expires_at = now + 30 days
  - Store in database

- [ ] 7.17 Handle retention = forever
  - expires_at = 0 (special value meaning "never expire")
  - Store in database

- [ ] 7.18 Handle retention = immediate
  - expires_at = now + 1 hour (give time to view)
  - Schedule for quick deletion

- [ ] 7.19 Handle custom retention (e.g., 7 days)
  - expires_at = now + N days
  - Store in database

---

### Cleanup Job Implementation

**Goal**: Automated job to delete expired transcripts.

**Reference**: 
- 📄 TDD Section 4 (Implementation Details - CleanupExpiredTranscripts)

#### Tasks:

- [ ] 7.20 Implement CleanupExpiredTranscripts() function
  ```go
  func (p *Plugin) CleanupExpiredTranscripts()
  ```

- [ ] 7.21 Run cleanup on schedule
  - Use ticker: `time.NewTicker(24 * time.Hour)`
  - Run once per day at 2 AM (off-peak)

- [ ] 7.22 Query expired transcripts
  ```sql
  SELECT id FROM meeting_transcripts
  WHERE expires_at > 0 AND expires_at < ?
  ```

- [ ] 7.23 Delete expired transcripts
  - Loop through results
  - Delete from meeting_transcripts
  - Also delete related meeting_summaries

- [ ] 7.24 Delete related recordings
  - Cascade delete meeting_recordings
  - Clean up orphaned records

- [ ] 7.25 Log cleanup results
  - Log number of transcripts deleted
  - Log database cleanup success

- [ ] 7.26 Handle cleanup errors
  - If delete fails, log error
  - Continue with other records
  - Retry on next run

- [ ] 7.27 Start cleanup job in OnActivate()
  - Launch cleanup goroutine
  - Store cleanup ticker in plugin struct

- [ ] 7.28 Stop cleanup job in OnDeactivate()
  - Stop ticker
  - Wait for cleanup to finish
  - Clean shutdown

---

### Testing

**Goal**: Verify settings and cleanup work correctly.

#### Tasks:

- [ ] 7.29 Test viewing settings
  - Run `/meeting-bot settings`
  - Verify default settings shown

- [ ] 7.30 Test updating retention to 7 days
  - Run `/meeting-bot settings retention 7`
  - Verify setting saved
  - Check database

- [ ] 7.31 Test setting retention to forever
  - Run `/meeting-bot settings retention forever`
  - Verify expires_at = 0 in database

- [ ] 7.32 Test setting retention to immediate
  - Run `/meeting-bot settings retention immediate`
  - Verify warning shown
  - Record meeting and verify transcript deleted quickly

- [ ] 7.33 Test retention applied to new transcripts
  - Change setting to 7 days
  - Record meeting
  - Check expires_at in database (should be now + 7 days)

- [ ] 7.34 Test cleanup job
  - Create transcript with expires_at = now - 1 day (already expired)
  - Manually trigger cleanup or wait for scheduled run
  - Verify transcript deleted

- [ ] 7.35 Test cleanup doesn't delete forever transcripts
  - Create transcript with expires_at = 0
  - Run cleanup
  - Verify transcript still exists

---

### PR #7 Wrap-Up

- [ ] 7.36 Document settings in README
  - List available settings
  - Explain retention options

- [ ] 7.37 Add privacy policy note
  - Explain data retention
  - User control over their data

- [ ] 7.38 Commit and create PR
  - Git add all files
  - Commit: "feat: User settings and automated cleanup"
  - Push to branch
  - Create PR

**PR #7 Complete** ✅

---

## PR #8: Testing Suite

**Branch**: `feature/meeting-bot-tests`  
**Description**: Comprehensive unit and integration tests for all components.

**Dependencies**: PRs #1-7 merged

**Success Criteria**:
- [ ] All core functions have unit tests
- [ ] Integration tests cover main flows
- [ ] Test coverage > 70%
- [ ] All tests pass
- [ ] CI/CD integration (if applicable)

---

### Unit Test Setup

**Goal**: Set up testing infrastructure.

#### Tasks:

- [ ] 8.1 Create test directory structure
  ```
  server/plugins/meeting-bot/
  ├── plugin_test.go
  ├── bot_test.go
  ├── commands_test.go
  ├── recorder_test.go
  ├── transcription_test.go
  ├── summarization_test.go
  ├── formatter_test.go
  ├── storage_test.go
  └── testutils/
      └── mocks.go
  ```

- [ ] 8.2 Create mock API plugin
  - Mock Mattermost API
  - Mock database calls
  - Mock HTTP clients

- [ ] 8.3 Create test fixtures
  - Sample audio files
  - Sample transcripts
  - Sample API responses

---

### Plugin Tests

**Goal**: Test plugin lifecycle.

#### Tasks:

- [ ] 8.4 Test OnActivate()
  - Verify bot user created
  - Verify commands registered
  - Verify database initialized

- [ ] 8.5 Test OnDeactivate()
  - Verify cleanup runs
  - Verify resources released

---

### Command Tests

**Goal**: Test all slash commands.

#### Tasks:

- [ ] 8.6 Test `/meeting-bot help`
  - Verify help text returned

- [ ] 8.7 Test `/meeting-bot start` with no call
  - Verify error message

- [ ] 8.8 Test `/meeting-bot start` with active call
  - Mock active call
  - Verify recording starts

- [ ] 8.9 Test `/meeting-bot settings`
  - Verify settings displayed

- [ ] 8.10 Test settings update
  - Verify setting saved to DB

---

### Transcription Tests

**Goal**: Test Whisper API integration.

#### Tasks:

- [ ] 8.11 Test TranscribeAudio() success
  - Mock Whisper API response
  - Verify transcript returned

- [ ] 8.12 Test TranscribeAudio() failure
  - Mock API error
  - Verify error handling

- [ ] 8.13 Test retry logic
  - Mock 2 failures, 1 success
  - Verify retries and final success

- [ ] 8.14 Test audio file deletion
  - Verify file deleted after transcription

---

### Summarization Tests

**Goal**: Test GPT-4o integration.

#### Tasks:

- [ ] 8.15 Test GenerateSummary() success
  - Mock GPT-4o response
  - Verify summary parsed correctly

- [ ] 8.16 Test JSON parsing
  - Mock response with markdown code blocks
  - Verify parser strips them

- [ ] 8.17 Test partial response handling
  - Mock response missing topics
  - Verify graceful degradation

---

### Formatter Tests

**Goal**: Test markdown formatting.

#### Tasks:

- [ ] 8.18 Test FormatResult()
  - Verify markdown structure
  - Verify all sections present

- [ ] 8.19 Test @mention parsing
  - Verify mentions formatted correctly

- [ ] 8.20 Test collapsible transcript
  - Verify <details> tags present

---

### Storage Tests

**Goal**: Test database operations.

#### Tasks:

- [ ] 8.21 Test storeTranscript()
  - Verify insertion
  - Verify expires_at calculated correctly

- [ ] 8.22 Test getUserSettings()
  - Verify default returned if no settings

- [ ] 8.23 Test cleanup job
  - Create expired transcript
  - Run cleanup
  - Verify deleted

---

### Integration Tests

**Goal**: Test end-to-end flows.

#### Tasks:

- [ ] 8.24 Test full recording flow
  - Start recording → end call → transcribe → summarize → post

- [ ] 8.25 Test concurrent recordings
  - Start 2 recordings
  - Verify both complete independently

- [ ] 8.26 Test error recovery
  - Fail transcription → retry → succeed

---

### PR #8 Wrap-Up

- [ ] 8.27 Run all tests
  - `go test ./...`
  - Verify all pass

- [ ] 8.28 Check test coverage
  - `go test -cover ./...`
  - Aim for > 70%

- [ ] 8.29 Add CI/CD configuration (optional)
  - GitHub Actions or similar
  - Run tests on every PR

- [ ] 8.30 Commit and create PR
  - Git add all test files
  - Commit: "test: Add comprehensive test suite"
  - Create PR

**PR #8 Complete** ✅

---

## Final Polish & Demo Preparation

**Goal**: Prepare for demo and final delivery.

**No separate PR - final cleanup on main branch**

---

### Documentation

#### Tasks:

- [ ] 9.1 Update main README
  - Complete feature list
  - Installation guide
  - Configuration guide
  - Usage examples

- [ ] 9.2 Create CHANGELOG.md
  - List all features
  - List all PRs
  - Version: 0.1.0

- [ ] 9.3 Create CONTRIBUTING.md (optional)
  - How to contribute
  - Development setup
  - Testing guidelines

- [ ] 9.4 Add screenshots
  - Recording in progress
  - Posted summary
  - Settings view

- [ ] 9.5 Add architecture diagram
  - High-level system overview
  - Component interactions

---

### Code Cleanup

#### Tasks:

- [ ] 9.6 Remove debug logging
  - Remove excessive logs
  - Keep important logs

- [ ] 9.7 Remove commented code
  - Clean up dead code

- [ ] 9.8 Format all code
  - Run `go fmt`
  - Run `goimports`

- [ ] 9.9 Add missing comments
  - Document exported functions
  - Add package comments

- [ ] 9.10 Run linter
  - `golangci-lint run`
  - Fix any issues

---

### Demo Preparation

#### Tasks:

- [ ] 9.11 Create demo script
  - Step-by-step demo flow
  - Talking points

- [ ] 9.12 Record test meeting
  - Clear audio
  - Mentions action items
  - 3-5 minutes long

- [ ] 9.13 Practice demo
  - Run through entire flow
  - Time it (should be < 5 minutes)

- [ ] 9.14 Prepare demo video
  - Record screen
  - Show: Command → Recording → Result
  - Add voiceover (optional)

- [ ] 9.15 Create demo slides (optional)
  - Problem statement
  - Solution overview
  - Architecture
  - Live demo
  - Technical challenges
  - Future enhancements

---

### Final Testing

#### Tasks:

- [ ] 9.16 Fresh install test
  - Uninstall plugin
  - Reinstall from scratch
  - Verify everything works

- [ ] 9.17 Test on clean database
  - Drop all plugin tables
  - Reinstall
  - Verify tables created

- [ ] 9.18 Test all commands
  - Help, start, stop, settings, retry
  - Verify all work

- [ ] 9.19 Test error scenarios
  - No API key
  - Invalid API key
  - No active call
  - All should fail gracefully

- [ ] 9.20 Performance test
  - Record 3 concurrent meetings
  - Verify no slowdown
  - Check memory usage

---

### Submission Preparation

#### Tasks:

- [ ] 9.21 Create final commit
  - "chore: Final polish for v0.1.0"

- [ ] 9.22 Create git tag
  - `git tag v0.1.0`

- [ ] 9.23 Write project reflection
  - What went well
  - What was challenging
  - What you learned
  - What you'd do differently

- [ ] 9.24 Update context files
  - Reflect final implementation
  - Note deviations from TDD

- [ ] 9.25 Prepare handoff document
  - Next steps / future work
  - Known issues / limitations
  - Deployment guide

---

## End of tasks-3.md

**All tasks complete!** 🎉

---

**Summary**:
- **tasks-1.md**: Investigation, Setup, Recording (PRs #1-2)
- **tasks-2.md**: Transcription, Summarization (PRs #3-4)
- **tasks-3.md**: Formatting, Error Handling, Settings, Testing, Demo (PRs #5-8 + Final Polish)

**Total**: 8 PRs + Final polish = Complete Meeting Bot Feature ✅

