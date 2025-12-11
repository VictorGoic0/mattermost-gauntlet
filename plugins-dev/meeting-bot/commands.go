package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

// registerCommand registers the /meeting-bot slash command
func (p *Plugin) registerCommand() error {
	command := &model.Command{
		Trigger:          "meeting-bot",
		AutoComplete:     true,
		AutoCompleteDesc: "Control the meeting assistant bot",
		AutoCompleteHint: "[start|stop|settings|help]",
		DisplayName:      "Meeting Assistant Bot",
		Description:      "AI-powered meeting transcription and summarization bot",
	}

	if err := p.API.RegisterCommand(command); err != nil {
		return fmt.Errorf("failed to register command: %w", err)
	}

	return nil
}

// ExecuteCommand handles slash command execution
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	// Parse command: /meeting-bot [subcommand]
	parts := strings.Fields(args.Command)
	if len(parts) < 1 {
		return p.handleHelp(args), nil
	}

	// Remove leading "/" from trigger
	trigger := strings.TrimPrefix(parts[0], "/")
	if trigger != "meeting-bot" {
		return nil, model.NewAppError("ExecuteCommand", "Invalid command trigger", nil, "", 400)
	}

	// Get subcommand (default to "help" if none provided)
	subcommand := "help"
	if len(parts) > 1 {
		subcommand = strings.ToLower(parts[1])
	}

	// Route to appropriate handler
	switch subcommand {
	case "start":
		return p.handleStart(args), nil
	case "stop":
		return p.handleStop(args), nil
	case "settings":
		return p.handleSettings(args), nil
	case "help", "":
		return p.handleHelp(args), nil
	default:
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Unknown subcommand: `%s`. Use `/meeting-bot help` for available commands.", subcommand),
		}, nil
	}
}

// handleHelp returns help message with available commands
func (p *Plugin) handleHelp(args *model.CommandArgs) *model.CommandResponse {
	helpText := "## Meeting Assistant Bot Commands\n\n" +
		"**Available Commands:**\n" +
		"- `/meeting-bot start` - Start recording the current call\n" +
		"- `/meeting-bot stop` - Stop recording the current call\n" +
		"- `/meeting-bot settings` - View or change your settings (coming soon)\n" +
		"- `/meeting-bot help` - Show this help message\n\n" +
		"**Usage:**\n" +
		"1. Join a voice call in this channel\n" +
		"2. Type `/meeting-bot start` to begin recording\n" +
		"3. The bot will join the call and record audio\n" +
		"4. When the call ends, you'll receive a transcript and summary\n\n" +
		"**Note:** Recording starts from when you run the command. Earlier conversation is not recorded."

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         helpText,
	}
}

// handleStart handles the start command (stub for now)
func (p *Plugin) handleStart(args *model.CommandArgs) *model.CommandResponse {
	// TODO: Implement full logic in PR #2
	// For now, just return a placeholder message
	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         "Start command received. Full implementation coming in PR #2.",
	}
}

// handleStop handles the stop command (stub for now)
func (p *Plugin) handleStop(args *model.CommandArgs) *model.CommandResponse {
	// TODO: Implement full logic in PR #2
	// For now, just return a placeholder message
	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         "Stop command received. Full implementation coming in PR #2.",
	}
}

// handleSettings handles the settings command (stub for now)
func (p *Plugin) handleSettings(args *model.CommandArgs) *model.CommandResponse {
	// TODO: Implement full logic in PR #7
	// For now, just return a placeholder message
	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         "Settings command coming soon. Full implementation in PR #7.",
	}
}

