package main

import (
	"fmt"

	"github.com/mattermost/mattermost/server/public/model"
)

const (
	botUsername    = "meeting-bot"
	botDisplayName = "Meeting Assistant"
	botDescription = "AI-powered meeting transcription bot"
	pluginID       = "com.mattermost.meeting-bot"
)

// ensureBotUser creates or ensures the bot user exists.
// Returns the bot user ID and any error encountered.
func (p *Plugin) ensureBotUser() (string, error) {
	botID, err := p.API.EnsureBotUser(&model.Bot{
		Username:    botUsername,
		DisplayName: botDisplayName,
		Description: botDescription,
		OwnerId:     pluginID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to ensure bot user: %w", err)
	}

	return botID, nil
}

