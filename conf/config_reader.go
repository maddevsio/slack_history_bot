package conf

import (
	"os"

	"github.com/gen1us2k/log"
	"github.com/urfave/cli"
)

// Version stores current service version
var (
	Version    string
	SlackToken string
	DBPath     string
	LogLevel   string
)

type Configuration struct {
	data *SlackBotConfig
	app  *cli.App
}

// NewConfigurator is constructor and creates a new copy of Configuration
func NewConfigurator() *Configuration {
	Version = "0.1dev"
	app := cli.NewApp()
	app.Name = "History bot for slack"
	app.Usage = "Saves and serves history"
	return &Configuration{
		data: &SlackBotConfig{},
		app:  app,
	}
}

func (c *Configuration) fillConfig() *SlackBotConfig {
	return &SlackBotConfig{
		SlackToken: SlackToken,
		DBPath:     DBPath,
	}
}

// Run is wrapper around cli.App
func (c *Configuration) Run() error {
	c.app.Before = func(ctx *cli.Context) error {
		log.SetLevel(log.MustParseLevel(LogLevel))
		return nil
	}
	c.app.Flags = c.setupFlags()
	return c.app.Run(os.Args)
}

// App is public method for Configuration.app
func (c *Configuration) App() *cli.App {
	return c.app
}

func (c *Configuration) setupFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "slack_token",
			Usage:       "Set slack token",
			EnvVar:      "SLACK_TOKEN",
			Destination: &SlackToken,
		},
		cli.StringFlag{
			Name:        "search_db",
			Value:       "history.bleve",
			Usage:       "Set database default",
			Destination: &DBPath,
			EnvVar:      "SEARCH_DB",
		},
		cli.StringFlag{
			Name:        "loglevel",
			Value:       "debug",
			Usage:       "set log level",
			Destination: &LogLevel,
			EnvVar:      "LOG_LEVEL",
		},
	}

}

// Get returns filled SlackBotConfig
func (c *Configuration) Get() *SlackBotConfig {
	c.data = c.fillConfig()
	return c.data
}
