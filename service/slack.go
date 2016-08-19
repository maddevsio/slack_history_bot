package service

import (
	"fmt"
	"strings"

	"github.com/gen1us2k/log"
	"github.com/nlopes/slack"
)

type SlackService struct {
	BaseService
	logger log.Logger
	sh     *SlackHistoryBot
	rtm    *slack.RTM
	me     string
}

func (ss *SlackService) Name() string {
	return "slack_worker"
}

func (ss *SlackService) Init(sh *SlackHistoryBot) error {
	ss.sh = sh
	ss.logger = log.NewLogger(ss.Name())
	api := slack.New(ss.sh.Config().SlackToken)
	//	logger := nlog.New(os.Stdout, "slack-bot: ", nlog.Lshortfile|nlog.LstdFlags)
	//	slack.SetLogger(logger)
	//	api.SetDebug(true)
	ss.rtm = api.NewRTM()
	return nil
}

func (ss *SlackService) Run() error {
	go ss.rtm.ManageConnection()
	for {
		if ss.me == "" {
			me := ss.rtm.GetInfo()
			if me != nil {
				ss.me = me.User.ID
				ss.logger.Info("I've found myself")
			}
		}
		select {
		case msg := <-ss.rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				//
			case *slack.ConnectedEvent:
				ss.rtm.SendMessage(ss.rtm.NewOutgoingMessage("Hello world", "C22UWDUQ3"))
			case *slack.MessageEvent:
				ss.logger.Infof("Message %s from channel %s from user %s at %s", ev.Msg.Text, ev.Channel, ev.Msg.User, ev.Msg.Timestamp)
				if ss.isToMe(ev.Msg.Text) {
					ss.logger.Info("I have a new message")
				}

			case *slack.PresenceChangeEvent:
				ss.logger.Infof("Presence Change: %v\n", ev)

			case *slack.LatencyReport:
				ss.logger.Infof("Current latency: %v\n", ev.Value)

			case *slack.RTMError:
				ss.logger.Infof("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				ss.logger.Infof("Invalid credentials")
				break

			default:
				//
				continue
			}
		}
	}
	return nil
}

func (ss *SlackService) isToMe(message string) bool {
	return strings.Contains(message, fmt.Sprintf("<@%s>", ss.me))
}
