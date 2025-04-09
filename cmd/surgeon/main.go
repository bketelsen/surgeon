package main

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/bketelsen/toolbox/cobra"
	"github.com/bketelsen/toolbox/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"go.uber.org/automaxprocs/maxprocs"

	goversion "github.com/bketelsen/toolbox/go-version"
	"github.com/spf13/viper"
	"github.com/thediveo/enumflag/v2"
)

var (
	// The global config object
	logLevel slog.Level
	appname  = "surgeon"

	version   = ""
	commit    = ""
	treeState = ""
	date      = ""
	builtBy   = ""
)

// LogLevelIDs Maps 3rd party enumeration values to their textual representations
var LogLevelIDs = map[slog.Level][]string{
	slog.LevelDebug: {"debug"},
	slog.LevelInfo:  {"info"},
	slog.LevelWarn:  {"warn"},
	slog.LevelError: {"error"},
}

func main() {
	cmd, config := NewRootCommand()
	cmd.AddCommand(NewInitCommand(config))
	cmd.AddCommand(NewCodemodCmd(config))
	cmd.AddCommand(NewManCommand(config))
	cmd.AddCommand(NewGendocsCommand(config))
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// ldflags
// Default: '-s -w -X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}} -X main.builtBy=goreleaser'.
var bversion = buildVersion(version, commit, date, builtBy, treeState)

// NewRootCommand creates a new root command for the application
func NewRootCommand() (*cobra.Command, *viper.Viper) {
	// this sets the config file location & file name to the current working directory
	// and the name of the package
	// this is the default location for the config file
	cwd, _ := os.Getwd()
	cfgFile := path.Join(cwd, fmt.Sprintf(".%s.yaml", appname))

	config := setupConfig()

	// Define our command
	rootCmd := &cobra.Command{
		Use:     "surgeon",
		Short:   "Surgically modify your forks",
		Version: bversion.String(),
		Long: `Surgeon is a tool to make surgical changes to forks of upstream repositories.

The surgeon command reads a configuration file in the current directory
named '.surgeon.yaml'.  This file contains the configuration for the
surgeon command.  The configuration file contains the upstream repository
URL, the directory containing the code modification files, and a list of
code modifications to apply to the forked repository.

The surgeon command will clone the upstream repository into a temporary directory,
then apply the code modifications to the cloned repository.  The contents of the
modified repository are copied to the current directory, overwriting any existing
files.

Important: modifications are applied in the order they are listed in the configuration,
and have a cumulative effect.  Be sure to verify your modifications before committing.`,

		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default slog logger to the cobra command
			slog.SetDefault(cmd.Logger)

			// bind the slog log level to our enumflag
			ll := config.GetInt("log-level")
			cmd.SetLogLevel(slog.Level(ll))

			// only prints if the log level is set to debug
			cmd.Logger.Debug("Debug logging enabled")

			// if the pflag has a value other than the default, then
			// reload the config file
			if cmd.Flags().Lookup("config-file").Changed {
				slog.Debug("Using config file from flag", "file", cfgFile)
				config.SetConfigFile(cfgFile)
				config.Set("config-file", cfgFile)
				return config.ReadInConfig()
			}
			// otherwise use the default config
			// created or loaded by the setupConfig function
			slog.Debug("Config Used", "file", config.ConfigFileUsed())

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			c, err := ReadConfig(config.GetString("config-file"))
			if err != nil {
				ui.Error("Specified config file not found", config.GetString("config-file"))
				return err
			}
			fmt.Println(c)
			cmd.Logger.Debug("config", "upstream", c.Upstream, "modsdir", c.ModsDir)
			project := NewPatient(c)
			return project.Operate()
		},
	}

	rootCmd.PersistentFlags().StringVarP(
		&cfgFile,
		"config-file",
		"c",
		cfgFile,
		``)
	_ = config.BindPFlag("config-file", rootCmd.PersistentFlags().Lookup("config-file"))
	rootCmd.PersistentFlags().Var(
		enumflag.New(&logLevel, "log", LogLevelIDs, enumflag.EnumCaseInsensitive),
		"log-level",
		"logging level [debug|info|warn|error]")
	return rootCmd, config
}

// https://www.asciiart.eu/text-to-ascii-art to make your own
// just make sure the font doesn't have backticks in the letters or
// it will break the string quoting
var asciiName = `
███████╗██╗   ██╗██████╗  ██████╗ ███████╗ ██████╗ ███╗   ██╗
██╔════╝██║   ██║██╔══██╗██╔════╝ ██╔════╝██╔═══██╗████╗  ██║
███████╗██║   ██║██████╔╝██║  ███╗█████╗  ██║   ██║██╔██╗ ██║
╚════██║██║   ██║██╔══██╗██║   ██║██╔══╝  ██║   ██║██║╚██╗██║
███████║╚██████╔╝██║  ██║╚██████╔╝███████╗╚██████╔╝██║ ╚████║
╚══════╝ ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚══════╝ ╚═════╝ ╚═╝  ╚═══╝
`

// buildVersion builds the version info for the application
func buildVersion(version, commit, date, builtBy, treeState string) goversion.Info {
	return goversion.GetVersionInfo(
		goversion.WithAppDetails(appname, "Surgically modify your forks.", "https://github.com/bketelsen/surgeon"),
		goversion.WithASCIIName(asciiName),
		func(i *goversion.Info) {
			if commit != "" {
				i.GitCommit = commit
			}
			if treeState != "" {
				i.GitTreeState = treeState
			}
			if date != "" {
				i.BuildDate = date
			}
			if version != "" {
				i.GitVersion = version
			}
			if builtBy != "" {
				i.BuiltBy = builtBy
			}
		},
	)
}

func init() {
	// enable colored output on github actions et al
	if os.Getenv("NOCOLOR") != "" {
		lipgloss.DefaultRenderer().SetColorProfile(termenv.Ascii)
	}
	// automatically set GOMAXPROCS to match available CPUs.
	// GOMAXPROCS will be used as the default value for the --parallelism flag.
	if _, err := maxprocs.Set(); err != nil {
		fmt.Println("failed to set GOMAXPROCS")
	}
}
