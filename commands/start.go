package commands

import (
	"context"
	"fmt"

	"github.com/amitizle/muffin/internal/logger"
	"github.com/amitizle/muffin/internal/scheduler"
	"github.com/amitizle/muffin/pkg/checks"
	"github.com/amitizle/muffin/pkg/notifiers"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start checks scheduler",
	Run:   startScheduler,
}

func init() {
	// startCmd.Flags().Bool("dry-run", false, "don't run the checks, just print that they were supposed to be running")
	rootCmd.AddCommand(startCmd)
}

func startScheduler(cmd *cobra.Command, args []string) {
	if err := initializeNotifiers(); err != nil {
		exitWithError(err)
	}

	s := scheduler.New()
	if err := initializeChecks(s); err != nil {
		exitWithError(err)
	}

	s.Start()
	select {}
}

func initializeChecks(s *scheduler.Scheduler) error {
	for _, cfgCheck := range cfg.Checks {
		checkLogger := log.With().Str("check_name", cfgCheck.Name).Str("check_type", cfgCheck.Type).Logger()
		checkLogger.Info().Msg("initializing check")
		check, err := checks.FromString(cfgCheck.Type)
		if err != nil {
			return err
		}

		ctx := context.Background()
		ctxWithLog := logger.StoreContext(ctx, checkLogger)

		if err := check.Initialize(ctxWithLog); err != nil {
			return err
		}

		if err := check.Configure(cfgCheck.Config); err != nil {
			return err
		}
		cfgCheck.Check = check
		err = s.NewTask(cfgCheck.Cron, func() {
			b, err := check.Run()
			if err != nil {
				log.Error().Err(err).Msg("failed check")
				for _, notifier := range cfg.Notifiers {
					notifier.Notifier.Notify(fmt.Sprintf("check %s error: %s", cfgCheck.Name, err))
				}
			}
			log.Info().Str("result", string(b)).Msg("check finished")
		})
		if err != nil {
			exitWithError(err)
		}
	}
	return nil
}

func initializeNotifiers() error {
	for _, cfgNotifier := range cfg.Notifiers {
		notifierLogger := log.With().Str("notifier", cfgNotifier.Name).Str("notifier_type", cfgNotifier.Type).Logger()
		notifierLogger.Info().Msg("initializing check")
		notifier, err := notifiers.FromString(cfgNotifier.Type)
		if err != nil {
			return err
		}

		ctx := context.Background()
		ctxWithLog := logger.StoreContext(ctx, notifierLogger)

		if err := notifier.Initialize(ctxWithLog); err != nil {
			return err
		}

		if err := notifier.Configure(cfgNotifier.Config); err != nil {
			return err
		}

		cfgNotifier.Notifier = notifier
	}
	return nil
}
