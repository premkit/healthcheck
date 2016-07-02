package commands

import (
	"fmt"

	"github.com/premkit/healthcheck/healthcheck"
	"github.com/premkit/healthcheck/log"
	"github.com/premkit/healthcheck/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultHTTPPort = 80

	defaultServiceFile            = ""
	defaultServiceFileContentType = "application/yaml"
	defaultDataFile               = "/data/healthcheck.db"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Create and run healthchecks on components",
	Long:  `Healthcheck is part of the Premkit package to deliver installable software`,
}

func init() {
	viper.SetEnvPrefix("healthcheck")
	viper.AutomaticEnv()

	serverCmd.Flags().Int("bind-http", defaultHTTPPort, "port on which the reverse proxy will bind and listen for http connections")
	serverCmd.Flags().String("data-file", defaultDataFile, "location of the database file")
	serverCmd.Flags().String("service-file", defaultServiceFile, "location of a service file to load")
	serverCmd.Flags().String("service-file-content-type", defaultServiceFileContentType, "content type of the service file to load")

	viper.BindPFlag("bind_http", serverCmd.Flags().Lookup("bind-http"))
	viper.BindPFlag("data_file", serverCmd.Flags().Lookup("data-file"))
	viper.BindPFlag("service_file", serverCmd.Flags().Lookup("service-file"))
	viper.BindPFlag("service_file_content_type", serverCmd.Flags().Lookup("service-file-content-type"))

	serverCmd.RunE = runServer
}

func runServer(cmd *cobra.Command, args []string) error {
	if err := InitializeConfig(serverCmd); err != nil {
		return err
	}

	config, err := buildConfig()
	if err != nil {
		return err
	}

	showAppliedSettings()

	server.Run(config)

	if viper.GetString("service_file") != "" {
		if _, err := healthcheck.ImportServiceFile(viper.GetString("service_file_content_type"), viper.GetString("service_file")); err != nil {
			return err
		}
	}

	<-make(chan int)

	return nil
}

func buildConfig() (*server.Config, error) {
	config := server.Config{
		HTTPPort: viper.GetInt("bind_http"),
	}

	return &config, nil
}

func showAppliedSettings() {
	var nonDefault []string

	if viper.GetInt("bind_http") != defaultHTTPPort {
		nonDefault = append(nonDefault, fmt.Sprintf("HTTP Bind Port set to %d", viper.GetInt("bind_http")))
	}
	if viper.GetString("data_file") != defaultDataFile {
		nonDefault = append(nonDefault, fmt.Sprintf("DataFile set to %s", viper.GetString("data_file")))
	}
	if viper.GetString("service_file") != defaultServiceFile {
		nonDefault = append(nonDefault, fmt.Sprintf("ServiceFile set to %s", viper.GetString("service_file")))
	}

	if len(nonDefault) == 0 {
		log.Infof("Using default settings")
		return
	}

	log.Infof("Overridden settings: ")
	for _, n := range nonDefault {
		log.Infof("\t%s", n)
	}
}
