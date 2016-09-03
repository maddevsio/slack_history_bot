package service

import (
	"fmt"
	"strings"

	"github.com/gen1us2k/log"
	"github.com/nlopes/slack"
)

type SlackService struct {
	BaseService
	logger   log.Logger
	sh       *SlackHistoryBot
	rtm      *slack.RTM
	search   *SearchService
	slackAPI *slack.Client
	me       string
}

func (ss *SlackService) Name() string {
	return "slack_worker"
}

func (ss *SlackService) Init(sh *SlackHistoryBot) error {
	ss.sh = sh
	ss.logger = log.NewLogger(ss.Name())
	api := slack.New(ss.sh.Config().SlackToken)
	ss.rtm = api.NewRTM()
	ss.search = ss.sh.SearchService()
	ss.slackAPI = api
	return nil
}

func (ss *SlackService) Run() error {
	go ss.rtm.ManageConnection()
	for {
		if ss.me == "" {
			me := ss.rtm.GetInfo()
			if me != nil {
				ss.me = me.User.ID
				ss.logger.Infof("I've found myself: %s", me.User.ID)
			}
		}
		select {
		case msg := <-ss.rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				ss.logger.Info("Bot connected!")
			case *slack.MessageEvent:
				if ss.isToMe(ev.Msg.Text) {
					ss.logger.Info("I have a new message")
					res, err := ss.search.Search(ss.cleanMessage(ev.Text), ev.Channel)
					if err != nil {
						ss.logger.Error(err)
					}
					message := fmt.Sprintf("+%v", res)
					ss.slackAPI.PostMessage(ev.Channel, message, slack.PostMessageParameters{})
					continue
				}
				ss.logger.Infof("Message %s from channel %s from user %s at %s", ev.Msg.Text, ev.Channel, ev.Msg.User, ev.Msg.Timestamp)
				ss.search.IndexMessage(IndexData{
					ID:        fmt.Sprintf("%s-%s", ev.Msg.User, ev.Msg.Timestamp),
					Username:  ev.Msg.User,
					Message:   ev.Msg.Text,
					Channel:   ev.Channel,
					Timestamp: ev.Timestamp,
				})

			case *slack.LatencyReport:
				ss.logger.Infof("Current latency: %v\n", ev.Value)

			case *slack.RTMError:
				ss.logger.Infof("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				ss.logger.Infof("Invalid credentials")
				break

			default:
				continue
			}
		}
	}
	return nil
}

func (ss *SlackService) isToMe(message string) bool {
	return strings.Contains(message, fmt.Sprintf("<@%s>", ss.me))
}

func (ss *SlackService) cleanMessage(message string) string {
	return strings.Replace(message, fmt.Sprintf("<@%s>", ss.me), "", -1)
}
