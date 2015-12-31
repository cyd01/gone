package config

import "fmt"

// Command represents a command that can be executed by the application.
type Command int

const (
	CommandHelp Command = iota
	CommandListen
	CommandExportTemplates
)

// String returns the string representation of the command, as it's to be used
// on the commandline.
func (c Command) String() string {
	switch c {
	case CommandHelp:
		return "help"
	case CommandListen:
		return "listen"
	case CommandExportTemplates:
		return "export-templates"
	}
	return ""
}

// Commands returns all valid command values.
func Commands() []Command {
	return []Command{CommandHelp, CommandListen, CommandExportTemplates}
}

// StringToCommand interprets the given string as a Command.
// It returns an error if the given string is no known command value.
func StringToCommand(s string) (Command, error) {
	for _, cmd := range Commands() {
		if s == cmd.String() {
			return cmd, nil
		}
	}
	return DefaultCommand, fmt.Errorf("invalid command: %s", s)
}

// Config is the configuration for the application.
// A Config instance might be created from different sources, like a
// configuration file or command line switches.
type Config interface {

	// Command is the command to be executed by the application.
	// This defaults to the DefaultCommand constant.
	Command() Command

	// BindAddress is the network address the application will listen on.
	// This defaults to the DefaultListenAddress constant.
	BindAddress() string

	// TemplatePath is the path to the directory containing custom templates.
	// This defaults to the empty string, meaning the static templates
	// delivered with the application are used.
	TemplatePath() string
}

const (
	DefaultCommand      = CommandListen
	DefaultBindAddress  = ":8080"
	DefaultTemplatePath = ""
)