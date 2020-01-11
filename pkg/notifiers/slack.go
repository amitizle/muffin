package notifiers

import (
	"context"
	"fmt"

	"github.com/amitizle/muffin/internal/logger"
	"github.com/mitchellh/mapstructure"
	"github.com/nlopes/slack"
	"github.com/rs/zerolog"
)

// SlackNotifier holds the configuration for an instance
// of slack notifier
type SlackNotifier struct {
	client *slack.Client
	config *slackNotifierConfig
	ctx    context.Context
	logger zerolog.Logger
}

type slackNotifierConfig struct {
	Token       string `mapstructure:"token"`
	ChannelName string `mapstructure:"channel_name"`
	channelID   string
}

func (notifier *SlackNotifier) Initialize(ctx context.Context) error {
	notifier.ctx = ctx
	lg, err := logger.GetContext(ctx)
	if err != nil {
		return err
	}
	notifier.logger = lg
	notifier.config = &slackNotifierConfig{}
	return nil
}

func (notifier *SlackNotifier) Configure(config map[string]interface{}) error {
	notifierConfig := &slackNotifierConfig{}
	if err := mapstructure.Decode(config, notifierConfig); err != nil {
		return err
	}
	notifier.client = slack.New(notifierConfig.Token)
	notifier.config = notifierConfig
	authResponse, err := notifier.client.AuthTest()
	if err != nil {
		return err
	}
	notifier.logger.Info().Msgf("authenticated as %s (team %s)", authResponse.User, authResponse.Team)

	channels, err := notifier.client.GetChannels(true)
	if err != nil {
		return err
	}

	var foundChannel bool
	for _, c := range channels {
		if c.NameNormalized == notifierConfig.ChannelName {
			notifierConfig.channelID = c.ID
			foundChannel = true
			break
		}
	}

	if !foundChannel {
		return fmt.Errorf("could not find channel ID for channel %s", notifierConfig.ChannelName)
	}
	return nil
}

func (notifier *SlackNotifier) Notify(msg string) error {
	_, _, err := notifier.client.PostMessage(notifier.config.channelID, slack.MsgOptionText(msg, false))
	if err != nil {
		return err
	}
	return nil
}
