# tasks-2.md

## Overview
This file covers Transcription (Whisper API) and Summarization (GPT-4o) implementation.

**PRs Covered**: #3-4
**Days**: Day 4
**Estimated Time**: ~8 hours total

---

## PR #3: Transcription Integration

**Branch**: `feature/meeting-bot-transcription`  
**Description**: Integrate OpenAI Whisper API for audio transcription with retry logic and error handling.

**Dependencies**: PR #2 merged (recording functionality working)

**Success Criteria**:
- [ ] Audio files successfully sent to Whisper API
- [ ] Transcripts returned and stored
- [ ] Retry logic works (tested with mock failures)
- [ ] Audio files deleted after successful transcription
- [ ] Rate limiting handled
- [ ] Error messages posted to users
- [ ] All transcription tests pass

---

### Whisper API Client Setup

**Goal**: Set up HTTP client for Whisper API requests.

**Reference**: 
- 📄 TDD Section 6 (API Integrations - Whisper API)
- 📄 TDD Section 4 (Implementation Details - transcription.go)

#### Tasks:

- [ ] 3.1 Create transcription.go file
  - Add to `server/plugins/meeting-bot/transcription.go`

- [ ] 3.2 Add HTTP client dependency
  - `go get github.com/go-resty/resty/v2`
  - Import in transcription.go

- [ ] 3.3 Define Whisper API constants
  ```go
  const (
      whisperAPIURL = "https://api.openai.com/v1/audio/transcriptions"
      whisperModel  = "whisper-1"
      maxRetries    = 3
  )
  ```

- [ ] 3.4 Define WhisperResponse struct
  ```go
  type WhisperResponse struct {
      Text string `json:"text"`
  }
  ```

- [ ] 3.5 Create HTTP client helper
  - Initialize resty client
  - Set timeout (5 minutes for large files)
  - Set retry policy (handled separately)

---

### Transcription Implementation

**Goal**: Send audio to Whisper and get transcript back.

**Reference**: 
- 📄 TDD Section 4 (Implementation Details - TranscribeAudio function)

#### Tasks:

- [ ] 3.6 Implement TranscribeAudio() function signature
  ```go
  func (p *Plugin) TranscribeAudio(audioPath string) (string, error)
  ```

- [ ] 3.7 Read audio file from disk
  - Use `ioutil.ReadFile(audioPath)`
  - Handle file not found error
  - Get file size for logging

- [ ] 3.8 Get OpenAI API key
  - Call `p.getOpenAIKey()`
  - Return error if key not configured
  - Never log the actual key

- [ ] 3.9 Prepare multipart form request
  - Create multipart writer
  - Add "file" field with audio data
  - Add "model" field: "whisper-1"
  - Add "response_format" field: "json"

- [ ] 3.10 Set request headers
  - Authorization: `Bearer {api_key}`
  - Content-Type: `multipart/form-data`

- [ ] 3.11 Send POST request to Whisper API
  - URL: `https://api.openai.com/v1/audio/transcriptions`
  - Include multipart body
  - Set timeout: 5 minutes

- [ ] 3.12 Handle HTTP errors
  - 400: Invalid audio format → return descriptive error
  - 401: Invalid API key → return "API key invalid"
  - 413: File too large → return "Audio file too large"
  - 429: Rate limit → return retryable error
  - 500: Server error → return retryable error
  - Other: Log full response, return generic error

- [ ] 3.13 Parse JSON response
  - Unmarshal response body
  - Extract "text" field
  - Handle malformed JSON

- [ ] 3.14 Validate transcript
  - Check transcript is not empty
  - Check length > 0
  - Log transcript length (chars/words)

- [ ] 3.15 Return transcript
  - Return transcript text
  - Return nil error on success

---

### Retry Logic Implementation

**Goal**: Retry failed transcriptions with exponential backoff.

**Reference**: 
- 📄 TDD Section 4 (Implementation Details - TranscribeWithRetry)
- 📄 TDD Section 9 (Error Handling - Retry Strategy)

#### Tasks:

- [ ] 3.16 Implement TranscribeWithRetry() function signature
  ```go
  func (p *Plugin) TranscribeWithRetry(audioPath string, channelID string) (string, error)
  ```

- [ ] 3.17 Implement retry loop
  - Max retries: 3
  - Loop from attempt 1 to maxRetries

- [ ] 3.18 Call TranscribeAudio() on each attempt
  - Pass audioPath
  - Capture error

- [ ] 3.19 Check if error is retryable
  - Rate limit (429): Yes, retry
  - Server error (500): Yes, retry
  - Timeout: Yes, retry
  - Invalid format (400): No, fail immediately
  - Invalid API key (401): No, fail immediately

- [ ] 3.20 Post retry notification to channel
  - Message: "⚠️ Transcription failed, retrying... (Attempt X/3)"
  - Only if not last attempt
  - Use bot user to post

- [ ] 3.21 Implement exponential backoff
  - Attempt 1: Wait 2 seconds
  - Attempt 2: Wait 4 seconds
  - Attempt 3: Wait 8 seconds
  - Use `time.Sleep()`

- [ ] 3.22 Return on success
  - If TranscribeAudio() succeeds, return immediately
  - Don't continue to next retry

- [ ] 3.23 Return final error after all retries
  - If all attempts fail, return last error
  - Include attempt count in error message

---

### Audio File Cleanup

**Goal**: Delete audio files immediately after transcription.

**Reference**: 
- 📄 TDD Section 7 (Security & Permissions - Data Privacy)

#### Tasks:

- [ ] 3.24 Implement deleteAudioFile() helper
  ```go
  func (p *Plugin) deleteAudioFile(audioPath string) error
  ```

- [ ] 3.25 Delete file from disk
  - Use `os.Remove(audioPath)`
  - Handle file not found (already deleted)
  - Log deletion success

- [ ] 3.26 Verify file deleted
  - Check file no longer exists
  - Log if deletion failed
  - Don't fail transcription if delete fails (just warn)

- [ ] 3.27 Call deleteAudioFile() after successful transcription
  - In TranscribeWithRetry(), after getting transcript
  - Delete regardless of success/failure
  - Log deletion result

- [ ] 3.28 Handle deletion errors gracefully
  - If delete fails, log warning
  - Post message to admin channel
  - Don't prevent transcript from being used

---

### Database Integration

**Goal**: Store transcripts in database.

**Reference**: 
- 📄 TDD Section 5 (Database Schema - meeting_transcripts)

#### Tasks:

- [ ] 3.29 Implement storeTranscript() function (basic version)
  ```go
  func (p *Plugin) storeTranscript(callID, recordingID, transcript string) error
  ```

- [ ] 3.30 Generate transcript ID
  - Use `model.NewId()`

- [ ] 3.31 Count words in transcript
  - Split by whitespace
  - Count tokens

- [ ] 3.32 Get user's retention setting
  - Query user_settings table
  - Default: 30 days
  - Calculate expires_at timestamp

- [ ] 3.33 Insert transcript into database
  ```sql
  INSERT INTO meeting_transcripts 
  (id, call_id, recording_id, transcript, word_count, created_at, expires_at)
  VALUES (?, ?, ?, ?, ?, ?, ?)
  ```

- [ ] 3.34 Handle insertion errors
  - Duplicate key: Log warning, don't fail
  - Other errors: Return error

- [ ] 3.35 Return transcript ID
  - Return newly created ID
  - Return error if failed

---

### Integration with Recording Flow

**Goal**: Call transcription after recording completes.

**Reference**: 
- 📄 TDD Section 3 (Data Flow - Processing Flow)

#### Tasks:

- [ ] 3.36 Update call_ended event handler
  - After SaveRecording() succeeds
  - Call processRecording()

- [ ] 3.37 Implement processRecording() function
  ```go
  func (p *Plugin) processRecording(recording *Recording) error
  ```

- [ ] 3.38 Update recording status in database
  - Set status = "processing"
  - Update updated_at timestamp

- [ ] 3.39 Post "Processing..." message
  - Message: "🤖 Processing your meeting... This may take 2 minutes."
  - Post to recording's channel

- [ ] 3.40 Get audio file path
  - Retrieve from recording.audio_file_path
  - Verify file exists

- [ ] 3.41 Call TranscribeWithRetry()
  - Pass audio path and channel ID
  - Capture transcript or error

- [ ] 3.42 Handle transcription success
  - Delete audio file
  - Store transcript in database
  - Update recording status = "transcribed"
  - Proceed to summarization (PR #4)

- [ ] 3.43 Handle transcription failure
  - Keep audio file (for manual retry)
  - Update recording status = "failed"
  - Post error message to channel
  - Store error details in database

---

### Error Messages

**Goal**: User-friendly error messages.

**Reference**: 
- 📄 TDD Section 9 (Error Handling - Error Messages)

#### Tasks:

- [ ] 3.44 Implement postErrorMessage() helper
  - Takes error and channel ID
  - Maps technical errors to user-friendly messages

- [ ] 3.45 Handle "API key not configured" error
  - Message: "❌ API key not configured. Admin: Configure in System Console > Plugins > Meeting Bot"
  - Post as bot user

- [ ] 3.46 Handle "API key invalid" error
  - Message: "❌ OpenAI API key is invalid. Admin: Check your API key in System Console."
  - Alert system admin

- [ ] 3.47 Handle "no speech detected" error
  - Check if transcript is empty or very short
  - Message: "⚠️ No speech detected in recording. This might be a technical issue. Please verify your audio setup."

- [ ] 3.48 Handle "file too large" error
  - Message: "❌ Audio file too large for transcription. Try recording shorter meetings."

- [ ] 3.49 Handle generic transcription failure
  - After 3 retries fail
  - Message: "❌ Transcription failed after 3 attempts. Audio has been saved. Admin can investigate or try `/meeting-bot retry {recording_id}`"

---

### Logging and Monitoring

**Goal**: Detailed logging for debugging.

#### Tasks:

- [ ] 3.50 Log transcription start
  - Log audio file path and size
  - Log recording ID and call ID

- [ ] 3.51 Log API request details
  - Log request URL (not full body)
  - Log file size being uploaded
  - Don't log API key

- [ ] 3.52 Log API response
  - Log HTTP status code
  - Log response time
  - Log transcript length

- [ ] 3.53 Log retry attempts
  - Log which attempt (1, 2, 3)
  - Log reason for retry
  - Log backoff duration

- [ ] 3.54 Log errors with context
  - Log full error message
  - Log call ID and channel ID
  - Log audio file path (for debugging)

- [ ] 3.55 Add metrics (optional)
  - Count successful transcriptions
  - Count failed transcriptions
  - Track average transcription time

---

### Testing

**Goal**: Verify transcription works end-to-end.

#### Tasks:

- [ ] 3.56 Test with valid audio file
  - Use recording from PR #2
  - Run transcription
  - Verify transcript returned

- [ ] 3.57 Verify transcript quality
  - Read transcript text
  - Check if it matches what was said
  - Note any mistakes (expected with Whisper)

- [ ] 3.58 Test with invalid audio file
  - Create empty file
  - Try to transcribe
  - Verify error handling

- [ ] 3.59 Test with missing audio file
  - Delete audio file before transcription
  - Verify error message

- [ ] 3.60 Test retry logic (mock failure)
  - Temporarily use invalid API key
  - Run transcription
  - Verify retry messages appear
  - Restore valid key and verify success

- [ ] 3.61 Test audio file deletion
  - Run full transcription
  - Verify audio file no longer exists in /tmp/

- [ ] 3.62 Test long audio (30+ min)
  - Record long meeting
  - Verify transcription completes
  - Check timeout doesn't trigger

- [ ] 3.63 Verify transcript in database
  - Check meeting_transcripts table
  - Verify all fields populated correctly
  - Check expires_at set properly

---

### PR #3 Wrap-Up

- [ ] 3.64 Update README
  - Document transcription process
  - Add troubleshooting for common errors

- [ ] 3.65 Add API key setup instructions
  - How to get OpenAI API key
  - Where to configure in Mattermost

- [ ] 3.66 Commit and create PR
  - Git add all files
  - Commit: "feat: Implement Whisper API transcription with retry logic"
  - Push to branch
  - Create PR with transcript example

**PR #3 Complete** ✅

---

## PR #4: Summarization Integration

**Branch**: `feature/meeting-bot-summarization`  
**Description**: Integrate GPT-4o API to extract action items, summary, and topics from transcripts.

**Dependencies**: PR #3 merged (transcription working)

**Success Criteria**:
- [ ] Transcripts successfully sent to GPT-4o
- [ ] Action items extracted correctly
- [ ] Summary is concise and accurate
- [ ] Topics identified
- [ ] Summaries stored in database
- [ ] JSON parsing robust (handles malformed responses)
- [ ] All summarization tests pass

---

### GPT-4o API Client Setup

**Goal**: Set up client for GPT-4o chat completions.

**Reference**: 
- 📄 TDD Section 6 (API Integrations - GPT-4o API)
- 📄 TDD Section 4 (Implementation Details - summarization.go)

#### Tasks:

- [ ] 4.1 Create summarization.go file
  - Add to `server/plugins/meeting-bot/summarization.go`

- [ ] 4.2 Define GPT-4o API constants
  ```go
  const (
      gpt4oAPIURL = "https://api.openai.com/v1/chat/completions"
      gpt4oModel  = "gpt-4o"
      temperature = 0.3
  )
  ```

- [ ] 4.3 Define Summary struct
  ```go
  type Summary struct {
      ActionItems []string `json:"action_items"`
      Summary     string   `json:"summary"`
      Topics      []string `json:"topics"`
  }
  ```

- [ ] 4.4 Define GPT-4o request/response structs
  ```go
  type ChatMessage struct {
      Role    string `json:"role"`
      Content string `json:"content"`
  }
  
  type ChatCompletionRequest struct {
      Model       string        `json:"model"`
      Messages    []ChatMessage `json:"messages"`
      Temperature float64       `json:"temperature"`
      MaxTokens   int           `json:"max_tokens"`
  }
  
  type ChatCompletionResponse struct {
      Choices []struct {
          Message ChatMessage `json:"message"`
      } `json:"choices"`
      Usage struct {
          TotalTokens int `json:"total_tokens"`
      } `json:"usage"`
  }
  ```

---

### Prompt Engineering

**Goal**: Create effective prompt for extracting meeting information.

**Reference**: 
- 📄 TDD Section 4 (Implementation Details - GenerateSummary)

#### Tasks:

- [ ] 4.5 Create system prompt
  ```go
  const systemPrompt = `You are a meeting assistant that extracts action items and summaries from meeting transcripts. Always return valid JSON.`
  ```

- [ ] 4.6 Create user prompt template
  ```go
  const userPromptTemplate = `Analyze this meeting transcript and extract:

1. Action Items: Specific tasks mentioned (with assignees if mentioned using @ notation)
2. Summary: 2-3 sentence overview of what was discussed
3. Topics: Key topics discussed (max 5)

Return ONLY valid JSON in this exact format:
{
  "action_items": ["@person to do X by Y", "Do Z (unassigned)"],
  "summary": "Brief summary here...",
  "topics": ["Topic 1", "Topic 2"]
}

Transcript:
%s`
  ```

- [ ] 4.7 Implement buildPrompt() helper
  - Takes transcript as input
  - Returns formatted prompt string
  - Uses fmt.Sprintf with template

- [ ] 4.8 Handle long transcripts
  - Check transcript length
  - If > 120K chars, truncate with warning
  - Log truncation for monitoring

---

### Summarization Implementation

**Goal**: Send transcript to GPT-4o and get structured summary.

**Reference**: 
- 📄 TDD Section 4 (Implementation Details - GenerateSummary function)

#### Tasks:

- [ ] 4.9 Implement GenerateSummary() function signature
  ```go
  func (p *Plugin) GenerateSummary(transcript string) (*Summary, error)
  ```

- [ ] 4.10 Build prompt from transcript
  - Call buildPrompt(transcript)
  - Get formatted prompt string

- [ ] 4.11 Create chat messages array
  - System message with systemPrompt
  - User message with built prompt

- [ ] 4.12 Build request body
  - Model: "gpt-4o"
  - Messages: [system, user]
  - Temperature: 0.3 (less random)
  - Max tokens: 2000 (enough for summary)

- [ ] 4.13 Get OpenAI API key
  - Call p.getOpenAIKey()
  - Return error if not configured

- [ ] 4.14 Send POST request to GPT-4o API
  - URL: chat completions endpoint
  - Headers: Authorization, Content-Type
  - Body: JSON request
  - Timeout: 60 seconds

- [ ] 4.15 Handle HTTP errors
  - 400: Bad request → log full error
  - 401: Invalid API key
  - 429: Rate limit → retry once after delay
  - 500: Server error → return error
  - Other: Log and return generic error

- [ ] 4.16 Parse response JSON
  - Unmarshal ChatCompletionResponse
  - Extract choices[0].message.content
  - Handle empty choices array

- [ ] 4.17 Extract content string
  - Get message content
  - Trim whitespace

- [ ] 4.18 Strip markdown code blocks
  - GPT-4o might wrap JSON in ```json ... ```
  - Remove opening: ```json\n
  - Remove closing: \n```
  - Trim result

- [ ] 4.19 Parse summary JSON
  - Unmarshal content into Summary struct
  - Handle JSON parse errors

- [ ] 4.20 Validate summary
  - Check action_items is array
  - Check summary is not empty
  - Check topics is array
  - If invalid, return error with details

- [ ] 4.21 Log token usage
  - Extract usage.total_tokens from response
  - Log for cost tracking
  - Calculate estimated cost

- [ ] 4.22 Return Summary struct
  - Return pointer to Summary
  - Return nil error on success

---

### Summary Storage

**Goal**: Store summaries in database.

**Reference**: 
- 📄 TDD Section 5 (Database Schema - meeting_summaries)

#### Tasks:

- [ ] 4.23 Update storeTranscript() to include summary
  - Rename to storeTranscriptAndSummary()
  - Take Summary as additional parameter

- [ ] 4.24 Insert summary into database
  ```sql
  INSERT INTO meeting_summaries
  (id, call_id, recording_id, transcript_id, summary, action_items, topics, created_at)
  VALUES (?, ?, ?, ?, ?, ?, ?, ?)
  ```

- [ ] 4.25 Serialize action_items to JSON
  - Use json.Marshal()
  - Store as JSONB/TEXT

- [ ] 4.26 Serialize topics to JSON
  - Use json.Marshal()
  - Store as JSONB/TEXT

- [ ] 4.27 Handle insertion errors
  - Log detailed error
  - Don't fail entire process if summary storage fails
  - Transcript still usable even if summary fails

- [ ] 4.28 Return summary ID
  - Return newly created summary ID
  - Use for linking in next steps

---

### Integration with Transcription Flow

**Goal**: Call summarization after transcription completes.

**Reference**: 
- 📄 TDD Section 3 (Data Flow - Processing Flow)

#### Tasks:

- [ ] 4.29 Update processRecording() function
  - After TranscribeWithRetry() succeeds
  - Call GenerateSummary()

- [ ] 4.30 Handle summarization success
  - Store transcript and summary together
  - Update recording status = "completed"
  - Proceed to formatting (PR #5)

- [ ] 4.31 Handle summarization failure
  - Store transcript anyway (it's still useful)
  - Update recording status = "transcribed_no_summary"
  - Post warning message
  - Don't fail entire process

- [ ] 4.32 Post intermediate status
  - After transcription: "🤖 Transcript ready. Generating summary..."
  - Gives user progress update

---

### Error Handling

**Goal**: Handle GPT-4o API errors gracefully.

**Reference**: 
- 📄 TDD Section 9 (Error Handling)

#### Tasks:

- [ ] 4.33 Handle JSON parsing errors
  - If GPT-4o returns invalid JSON
  - Log the raw response
  - Try to extract partial data
  - Fall back to empty summary if needed

- [ ] 4.34 Handle malformed JSON structure
  - If JSON is valid but missing fields
  - Use defaults for missing fields
  - Log warning

- [ ] 4.35 Handle rate limits
  - If 429 error, wait 10 seconds
  - Retry once
  - If still fails, return error

- [ ] 4.36 Handle timeout errors
  - If request takes > 60 seconds
  - Post message: "Summary generation is taking longer than expected..."
  - Retry with shorter transcript (truncated)

- [ ] 4.37 Handle empty action items
  - If GPT-4o returns no action items
  - Accept it (meeting might have no actions)
  - Don't treat as error

---

### Logging

**Goal**: Log summarization process for debugging.

#### Tasks:

- [ ] 4.38 Log summarization start
  - Log transcript length
  - Log recording ID

- [ ] 4.39 Log API request
  - Log prompt length
  - Log model being used
  - Don't log full transcript (too long)

- [ ] 4.40 Log API response
  - Log HTTP status
  - Log response time
  - Log token usage and cost

- [ ] 4.41 Log parsed summary
  - Log number of action items
  - Log summary length (chars)
  - Log number of topics

- [ ] 4.42 Log errors with context
  - Log full error message
  - Log recording ID
  - Log partial response if available

---

### Testing

**Goal**: Verify summarization works correctly.

#### Tasks:

- [ ] 4.43 Test with real transcript
  - Use transcript from PR #3
  - Run summarization
  - Verify summary returned

- [ ] 4.44 Verify action items extracted
  - Read action items list
  - Check if they match what was discussed
  - Verify assignees captured if mentioned

- [ ] 4.45 Verify summary quality
  - Read 2-3 sentence summary
  - Check if it captures main points
  - Verify it's concise

- [ ] 4.46 Verify topics identified
  - Check topics list
  - Verify topics are relevant
  - Check max 5 topics returned

- [ ] 4.47 Test with transcript with no action items
  - Create transcript with just discussion
  - Verify empty action_items array acceptable

- [ ] 4.48 Test with very long transcript
  - Create 2-hour meeting transcript
  - Verify truncation if needed
  - Verify still gets summary

- [ ] 4.49 Test JSON parsing edge cases
  - Mock response with markdown code blocks
  - Mock response with extra whitespace
  - Verify parser handles both

- [ ] 4.50 Verify summary in database
  - Check meeting_summaries table
  - Verify JSONB fields stored correctly
  - Query action_items and topics

---

### PR #4 Wrap-Up

- [ ] 4.51 Update README
  - Document summarization process
  - Add examples of summaries

- [ ] 4.52 Document prompt engineering
  - Save prompts in documentation
  - Note that prompts can be customized

- [ ] 4.53 Commit and create PR
  - Git add all files
  - Commit: "feat: Implement GPT-4o summarization with action items extraction"
  - Push to branch
  - Create PR with example summary

**PR #4 Complete** ✅
