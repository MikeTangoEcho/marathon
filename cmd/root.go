package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Logging
var logLevel string

func SetupLog() {
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		// TODO Better handle log level validation
		panic(err)
	}
	log.SetLevel(level)
}

// Config
//const defaultCfgFilename string = "twitch-marathon"

//var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "marathon",
	Short: "Chat driven live streaming.",
	Long: `Run a bot that will read chats for commands.

	!play *file*	Stream the content of the file according to the chosen tool.
	!shutdown	Close the stream and exists the application.
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		SetupLog()
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func SetVersion(version string) {
	rootCmd.Version = version
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/"+defaultCfgFilename+".yaml)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", log.WarnLevel.String(), "Log level")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Binds
	viper.BindPFlag("log-level", rootCmd.Flags().Lookup("log-level"))
}
