package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	logLevel  string
	logger    *zap.Logger
	rootViper = viper.New()
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rmdup",
	Short: "the template of cobra",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application`,

	PreRunE: preRunE,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		if h, _ := cmd.Flags().GetBool("help"); h {
			return cmd.Help()
		}
		return run(cmd, args)
	},
	SilenceErrors: true,
	SilenceUsage:  true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
	}
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().BoolP("verbose", "V", false, "verbose")

	rootCmd.PersistentFlags().StringVar(&logLevel, "log.level", "debug", "level of logger")
}

func preRunE(cmd *cobra.Command, args []string) (err error) {

	// use flag log.level
	if strings.ToLower(logLevel) == "debug" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return err
	}

	// Viper uses the following precedence order. Each item takes precedence over the item below it:
	//
	// explicit call to Set
	// flag
	// env
	// config
	// key/value store
	// default
	//
	// Viper configuration keys are case insensitive.

	return nil
}
