package cmd

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/davepgreene/propsd/http"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/davepgreene/propsd/config"
)

var cfgFile string
var verbose bool

// PropsdCmd represents the base command when called without any subcommands
var PropsdCmd = &cobra.Command{
	Use:   "propsd",
	Short: "Dynamic property management at scale",
	Long: `Propsd does dynamic property management at scale, across
	thousands of servers and changes from hundreds of developers, leveraging
	Amazon S3 to deliver properties and Consul to handle service discovery.
	Composable layering lets you set properties for an organization, a single
	server, and everything in between. Plus, flat file storage makes backups
	and audits a breeze.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := initializeConfig()
		initializeLog()
		if err != nil {
			return err
		}

		return boot()
	},
}

func boot() error {
	router := http.Handler()
	log.Error(router)
	return router
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the PropsdCmd.
func Execute() {
	if err := PropsdCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	PropsdCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
	PropsdCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose level logging")
	validConfigFilenames := []string{"json"}
	PropsdCmd.PersistentFlags().SetAnnotation("config", cobra.BashCompFilenameExt, validConfigFilenames)
}

func initializeLog() {
	log.RegisterExitHandler(func() {
		log.Info("Shutting down")
	})

	// Set logging options based on config
	if lvl, err := log.ParseLevel(viper.GetString("log.level")); err == nil {
		log.SetLevel(lvl)
	} else {
		log.Info("Unable to parse log level in settings. Defaulting to INFO")
	}

	// If using verbose mode, log at debug level
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	if viper.GetBool("log.json") {
		log.SetFormatter(&log.JSONFormatter{})
	}

	if cfgFile != "" {
		log.WithFields(log.Fields{
			"file": viper.ConfigFileUsed(),
		}).Info("Loaded config file")
	}

}

func initializeConfig(subCmdVs ...*cobra.Command) error {
	config.Defaults()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	viper.AutomaticEnv() // read in environment variables that match

	return nil
}
