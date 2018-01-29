package godiscordbot

import "github.com/bwmarrin/discordgo"

// Module represents any type of plugin that we can run
// we designi plugins to isolate features and allow them
// to register events that each one needs
type Module interface {
	Register(*discordgo.Session) error
}

// Modules is a list of all initalized modules at runtime
var modules []Module

// Initalize will bootup every module
func Initalize(s *discordgo.Session) error {
	for _, m := range modules {
		if err := m.Register(s); err != nil {
			return err
		}
	}

	return nil
}

func AddModule(m Module) {
	modules = append(modules, m)
}
