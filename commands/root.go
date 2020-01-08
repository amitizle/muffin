package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/amitizle/muffin/internal/config"
	"github.com/amitizle/muffin/internal/logger"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	cfg *config.Config

	rootCmd = &cobra.Command{
		Use:   "muffin",
		Short: "Muffin is a simple use application to monitor network services",
	}
)

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitWithError(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initLog)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.muffin.yaml)")
}

func initLog() {
	if err := logger.Init(cfg.Log.Level); err != nil {
		exitWithError(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfg = config.New()
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			exitWithError(err)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".muffin")
	}

	viper.SetEnvPrefix("muffin")
	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(cfg); err != nil {
		exitWithError(err)
	}
}

func exitWithError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
