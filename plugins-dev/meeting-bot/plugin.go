package main

import (
	"fmt"

	"github.com/mattermost/mattermost/server/public/plugin"
)

type Plugin struct {
	plugin.MattermostPlugin
	botUserID string
}

func (p *Plugin) OnActivate() error {
	p.API.LogInfo("Meeting Assistant Bot plugin activated")

	// Create bot user
	botID, err := p.ensureBotUser()
	if err != nil {
		return fmt.Errorf("failed to create bot user: %w", err)
	}
	p.botUserID = botID
	p.API.LogInfo("Bot user created", "bot_user_id", botID)

	// Register slash command
	if err := p.registerCommand(); err != nil {
		return fmt.Errorf("failed to register command: %w", err)
	}
	p.API.LogInfo("Slash command registered successfully")

	// TODO: Initialize database
	return nil
}

func (p *Plugin) OnDeactivate() error {
	p.API.LogInfo("Meeting Assistant Bot plugin deactivated")
	// TODO: Clean up resources
	return nil
}

func main() {
	plugin.ClientMain(&Plugin{})
}

