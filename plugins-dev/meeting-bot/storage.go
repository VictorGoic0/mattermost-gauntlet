package main

import (
	"database/sql"
	"fmt"

	"github.com/mattermost/mattermost/server/public/pluginapi"
)

// initDatabase creates all required database tables if they don't exist
func (p *Plugin) initDatabase() error {
	p.API.LogInfo("Initializing database tables...")

	// Use pluginapi client to get StoreService
	if p.Driver == nil {
		return fmt.Errorf("plugin driver is nil - database access not available")
	}

	client := pluginapi.NewClient(p.API, p.Driver)
	store := client.Store
	defer store.Close()

	db, err := store.GetMasterDB()
	if err != nil {
		p.API.LogError("Failed to get database connection", "error", err.Error())
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// Create all tables
	if err := p.createMeetingRecordingsTable(db); err != nil {
		p.API.LogError("Failed to create meeting_recordings table", "error", err.Error())
		return fmt.Errorf("failed to create meeting_recordings table: %w", err)
	}

	if err := p.createMeetingTranscriptsTable(db); err != nil {
		p.API.LogError("Failed to create meeting_transcripts table", "error", err.Error())
		return fmt.Errorf("failed to create meeting_transcripts table: %w", err)
	}

	if err := p.createMeetingSummariesTable(db); err != nil {
		p.API.LogError("Failed to create meeting_summaries table", "error", err.Error())
		return fmt.Errorf("failed to create meeting_summaries table: %w", err)
	}

	if err := p.createPluginConfigTable(db); err != nil {
		p.API.LogError("Failed to create plugin_config table", "error", err.Error())
		return fmt.Errorf("failed to create plugin_config table: %w", err)
	}

	if err := p.createUserSettingsTable(db); err != nil {
		p.API.LogError("Failed to create user_settings table", "error", err.Error())
		return fmt.Errorf("failed to create user_settings table: %w", err)
	}

	p.API.LogInfo("Database tables initialized successfully")
	return nil
}

// createMeetingRecordingsTable creates the meeting_recordings table
func (p *Plugin) createMeetingRecordingsTable(db *sql.DB) error {
	query := `
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
	`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	// Create indexes
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_meeting_recordings_call_id ON meeting_recordings(call_id);",
		"CREATE INDEX IF NOT EXISTS idx_meeting_recordings_status ON meeting_recordings(status);",
		"CREATE INDEX IF NOT EXISTS idx_meeting_recordings_started_by ON meeting_recordings(started_by);",
	}

	for _, indexQuery := range indexes {
		if _, err := db.Exec(indexQuery); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}

// createMeetingTranscriptsTable creates the meeting_transcripts table
func (p *Plugin) createMeetingTranscriptsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS meeting_transcripts (
			id VARCHAR(26) PRIMARY KEY,
			call_id VARCHAR(26) NOT NULL,
			recording_id VARCHAR(26) NOT NULL,
			transcript TEXT NOT NULL,
			word_count INTEGER,
			created_at BIGINT NOT NULL,
			expires_at BIGINT
		);
	`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	// Create indexes
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_meeting_transcripts_call_id ON meeting_transcripts(call_id);",
		"CREATE INDEX IF NOT EXISTS idx_meeting_transcripts_expires_at ON meeting_transcripts(expires_at);",
		"CREATE INDEX IF NOT EXISTS idx_meeting_transcripts_created_at ON meeting_transcripts(created_at DESC);",
	}

	for _, indexQuery := range indexes {
		if _, err := db.Exec(indexQuery); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}

// createMeetingSummariesTable creates the meeting_summaries table
func (p *Plugin) createMeetingSummariesTable(db *sql.DB) error {
	query := `
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
	`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	// Create indexes
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_meeting_summaries_call_id ON meeting_summaries(call_id);",
		"CREATE INDEX IF NOT EXISTS idx_meeting_summaries_recording_id ON meeting_summaries(recording_id);",
	}

	for _, indexQuery := range indexes {
		if _, err := db.Exec(indexQuery); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}

// createPluginConfigTable creates the plugin_config table
func (p *Plugin) createPluginConfigTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS plugin_config (
			key VARCHAR(255) PRIMARY KEY,
			value TEXT NOT NULL,
			encrypted BOOLEAN DEFAULT FALSE,
			updated_at BIGINT NOT NULL
		);
	`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}

// createUserSettingsTable creates the user_settings table
func (p *Plugin) createUserSettingsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS user_settings (
			user_id VARCHAR(26) PRIMARY KEY,
			retention_days INTEGER DEFAULT 30,
			auto_join BOOLEAN DEFAULT FALSE,
			created_at BIGINT NOT NULL,
			updated_at BIGINT NOT NULL
		);
	`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}

