package service

import (
	"fmt"
	"sync"

	"github.com/gen1us2k/log"

	"github.com/maddevsio/slackhistorybot/conf"
)

// SlackHistoryBot is main struct of daemon
// it stores all services that used by
type SlackHistoryBot struct {
	config *conf.SlackBotConfig

	services  map[string]Service
	waitGroup sync.WaitGroup

	logger log.Logger
}

// NewSlackHistoryBot creates and returns new SlackHistoryBotInstance
func NewSlackHistoryBot(config *conf.SlackBotConfig) *SlackHistoryBot {
	pb := new(SlackHistoryBot)
	pb.config = config
	pb.logger = log.NewLogger("slack_history_bot")
	pb.services = make(map[string]Service)
	pb.AddService(&SearchService{})
	pb.AddService(&SlackService{})
	return pb
}

// Start starts all services in separate goroutine
func (pb *SlackHistoryBot) Start() error {
	pb.logger.Info("Starting bot service")
	for _, service := range pb.services {
		pb.logger.Infof("Initializing: %s\n", service.Name())
		if err := service.Init(pb); err != nil {
			return fmt.Errorf("initialization of %q finished with error: %v", service.Name(), err)
		}
		pb.waitGroup.Add(1)

		go func(srv Service) {
			defer pb.waitGroup.Done()
			pb.logger.Infof("running %q service\n", srv.Name())
			if err := srv.Run(); err != nil {
				pb.logger.Errorf("error on run %q service, %v", srv.Name(), err)
			}
		}(service)
	}
	return nil
}

// AddService adds service into SlackHistoryBot.services map
func (pb *SlackHistoryBot) AddService(srv Service) {
	pb.services[srv.Name()] = srv

}

// Config returns current instance of SlackHistoryBotConfig
func (pb *SlackHistoryBot) Config() conf.SlackBotConfig {
	return *pb.config
}

// Stop stops all services running
func (pb *SlackHistoryBot) Stop() {
	pb.logger.Info("Worker is stopping...")
	for _, service := range pb.services {
		service.Stop()
	}
}

// WaitStop blocks main thread and waits when all goroutines will be stopped
func (pb *SlackHistoryBot) WaitStop() {
	pb.waitGroup.Wait()
}

func (pb *SlackHistoryBot) SearchService() *SearchService {
	service, ok := pb.services["search_service"]
	if !ok {
		pb.logger.Error("search service not found")
	}
	return service.(*SearchService)
}
