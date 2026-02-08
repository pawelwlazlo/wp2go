package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configPath      string
	hostName        string
	configErr       error
	forceConfigured bool
)

func initConfig() {
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("WP2GO")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if configPath != "" {
		configPath = expandPath(configPath)
		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			configErr = fmt.Errorf("read config: %w", err)
			return
		}
	} else {
		viper.SetConfigName("wp2go")
		viper.AddConfigPath(".")
		if home, err := os.UserHomeDir(); err == nil {
			viper.AddConfigPath(filepath.Join(home, ".config", "wp2go"))
		}
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				configErr = fmt.Errorf("read config: %w", err)
				return
			}
		}
	}

	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		configErr = fmt.Errorf("bind flags: %w", err)
	}
}

func applyConfig(cmd *cobra.Command) error {
	if configErr != nil {
		return configErr
	}

	if hostName == "" {
		if viper.IsSet("host") {
			hostName = viper.GetString("host")
		} else if auto, ok := autoDetectHost(); ok {
			hostName = auto
		}
	}

	if hostName != "" && !hostExists(hostName) {
		return fmt.Errorf("host %q not found in config", hostName)
	}

	if !cmd.Flags().Changed("archive") {
		if v, ok := getConfigString(hostName, "archive"); ok {
			archivePath = v
		}
	}
	if !cmd.Flags().Changed("db") {
		if v, ok := getConfigString(hostName, "db"); ok {
			dbPath = v
		}
	}
	if !cmd.Flags().Changed("sitename") {
		if v, ok := getConfigString(hostName, "sitename"); ok {
			sitename = v
		}
	}
	if !cmd.Flags().Changed("output-path") {
		if v, ok := getConfigString(hostName, "output_path"); ok {
			outputPath = v
		}
	}
	if !cmd.Flags().Changed("force") {
		if v, ok := getConfigBool(hostName, "force"); ok {
			force = v
			forceConfigured = true
		}
	}

	if sitename == "" && hostName != "" && !cmd.Flags().Changed("sitename") {
		sitename = hostName
	}

	return nil
}

func autoDetectHost() (string, bool) {
	hosts := viper.GetStringMap("hosts")
	if len(hosts) != 1 {
		return "", false
	}
	for name := range hosts {
		return name, true
	}
	return "", false
}

func hostExists(name string) bool {
	hosts := viper.GetStringMap("hosts")
	if len(hosts) == 0 {
		return false
	}
	_, ok := hosts[name]
	return ok
}

func getConfigString(host, key string) (string, bool) {
	if host != "" {
		hostKey := fmt.Sprintf("hosts.%s.%s", host, key)
		if viper.IsSet(hostKey) {
			return viper.GetString(hostKey), true
		}
	}
	defaultKey := fmt.Sprintf("defaults.%s", key)
	if viper.IsSet(defaultKey) {
		return viper.GetString(defaultKey), true
	}
	return "", false
}

func getConfigBool(host, key string) (bool, bool) {
	if host != "" {
		hostKey := fmt.Sprintf("hosts.%s.%s", host, key)
		if viper.IsSet(hostKey) {
			return viper.GetBool(hostKey), true
		}
	}
	defaultKey := fmt.Sprintf("defaults.%s", key)
	if viper.IsSet(defaultKey) {
		return viper.GetBool(defaultKey), true
	}
	return false, false
}
