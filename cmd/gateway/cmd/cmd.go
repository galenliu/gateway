// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"fmt"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	optionNameDataDir            = "data-dir"
	optionNameMediaDir           = "media-dir"
	optionNameLogDir             = "log-dir"
	optionNameAttachAddonDirs    = "attach-addons-dir"
	optionNameVerbosity          = "verbosity"
	optionNameDBRemoveBeforeOpen = "db-remove-before-open"
	optionNameAPIPort            = "api-port"
	optionNameIpcPort            = "ipc-port"
	optionNameRpcPort            = "rpc-port"
	optionNameHttpPort           = "http-port"
	optionNameHttpsPort          = "https-port"
	optionNameAddonUrls          = "addon-urls"
	optionLogRotateDays          = "log-rotate-days"
	optionHomeKitPin             = "hk-pin"
	optionHomeKitEnable          = "homekit-enable"
)

func init() {
	cobra.EnableCommandSorting = false
}

type command struct {
	root    *cobra.Command
	config  *viper.Viper
	cfgFile string
	homeDir string
}

type option func(*command)

func newCommand(opts ...option) (c *command, err error) {
	c = &command{
		root: &cobra.Command{
			Use:           "gateway",
			Short:         "WebThings Gateway",
			SilenceErrors: true,
			SilenceUsage:  true,
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				return c.initConfig()
			},
		},
	}

	for _, o := range opts {
		o(c)
	}

	// Find home directory.
	if err := c.setHomeDir(); err != nil {
		return nil, err
	}

	c.initGlobalFlags()

	if err := c.initStartCmd(); err != nil {
		return nil, err
	}

	if err := c.initInitCmd(); err != nil {
		return nil, err
	}

	c.initVersionCmd()

	if err := c.initConfiguratorOptionsCmd(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *command) Execute() (err error) {
	return c.root.Execute()
}

// Execute parses command line arguments and runs appropriate functions.
func Execute() (err error) {
	c, err := newCommand()
	if err != nil {
		return err
	}
	return c.Execute()
}

func (c *command) initGlobalFlags() {
	globalFlags := c.root.PersistentFlags()
	globalFlags.StringVar(&c.cfgFile, "config", "", "config file (default is $HOME/.gateway.yaml)")
}

func (c *command) initConfig() (err error) {
	config := viper.New()
	configName := ".gateway"
	if c.cfgFile != "" {
		// Use config file from the flag.
		config.SetConfigFile(c.cfgFile)
	} else {
		// Search config in home directory with name ".gateway" (without extension).
		config.AddConfigPath(c.homeDir)
		config.SetConfigName(configName)
	}

	// Environment
	config.SetEnvPrefix("gateway")
	config.AutomaticEnv() // read in environment variables that match
	config.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if c.homeDir != "" && c.cfgFile == "" {
		c.cfgFile = filepath.Join(c.homeDir, configName+".yaml")
	}

	// If a config file is found, read it in.
	if err := config.ReadInConfig(); err != nil {
		var e viper.ConfigFileNotFoundError
		if !errors.As(err, &e) {
			return err
		}
	}
	c.config = config
	return nil
}

func (c *command) setHomeDir() (err error) {
	if c.homeDir != "" {
		return
	}
	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	c.homeDir = dir
	return nil
}

func (c *command) setAllFlags(cmd *cobra.Command) {
	cmd.Flags().String(optionNameDataDir, filepath.Join(c.homeDir, ".gateway"), "data directory")

	dataDir, _ := cmd.Flags().GetString(optionNameDataDir)
	cmd.Flags().String(optionNameMediaDir, filepath.Join(dataDir, "media"), "media directory")

	cmd.Flags().String(optionNameLogDir, filepath.Join(dataDir, "log"), "media directory")
	cmd.Flags().String(optionNameAttachAddonDirs, "", "add-ons directory")
	cmd.Flags().Bool(optionNameDBRemoveBeforeOpen, false, "remove db before open")
	cmd.Flags().String(optionNameVerbosity, "info", "log verbosity level 0=silent, 1=error, 2=warn, 3=info, 4=debug, 5=trace")
	cmd.Flags().StringArray(optionNameAddonUrls, []string{"https://api.webthings.io:8443/addons"}, "addon urls")

	cmd.Flags().Int(optionNameAPIPort, 9090, "HTTP API listen address")
	cmd.Flags().Int(optionNameHttpPort, 9090, "http port")
	cmd.Flags().Int(optionNameHttpsPort, 4443, "https port")

	cmd.Flags().Int(optionNameIpcPort, 9500, "ipc port")
	cmd.Flags().Int(optionNameRpcPort, 9600, "rpc port")

	cmd.Flags().Int(optionLogRotateDays, 7, "log rotate days")
	cmd.Flags().String(optionHomeKitPin, "12344321", "homekit pin")
	cmd.Flags().Bool(optionHomeKitEnable, true, "homekit enable")
}

func newLogger(cmd *cobra.Command, verbosity string) (logging.Logger, error) {
	var logger logging.Logger
	switch verbosity {
	case "0", "silent":
		logger = logging.New(ioutil.Discard, 0)
	case "1", "error":
		logger = logging.New(cmd.OutOrStdout(), logrus.ErrorLevel)
	case "2", "warn":
		logger = logging.New(cmd.OutOrStdout(), logrus.WarnLevel)
	case "3", "info":
		logger = logging.New(cmd.OutOrStdout(), logrus.InfoLevel)
	case "4", "debug":
		logger = logging.New(cmd.OutOrStdout(), logrus.DebugLevel)
	case "5", "trace":
		logger = logging.New(cmd.OutOrStdout(), logrus.TraceLevel)
	default:
		return nil, fmt.Errorf("unknown verbosity level %q", verbosity)
	}
	return logger, nil
}
