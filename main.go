package main

import (
	"errors"
	"fmt"
	"github.com/Go-routine-4595/how-oem-view/controllers/api"
	"github.com/Go-routine-4595/how-oem-view/controllers/gateway/datasource"
	"github.com/Go-routine-4595/how-oem-view/model"
	"github.com/rs/zerolog"
	_ "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	datasource.EventHubConfig `yaml:"EventHubConfig"`
}

func main() {
	var (
		service    model.IService
		apiSrv     api.ApiController
		data       *datasource.DataSource
		mlog       model.IService
		conf       Config
		configFile string
		dataFile   string
		rootCmd    *cobra.Command
		flageh     bool
		flagfile   bool
	)

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	rootCmd = &cobra.Command{
		Use:   "oem-sim-viewer-server",
		Short: "A simple CLI app to view OEM alarms, it accepts a config file (default config.yaml)",
		Long:  "A simple CLI app to view OEM alarms, it read from a jsonl OEM file or an Event Hub broker in this case it requires a config file (default config.yaml)",
	}
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file (default is config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&dataFile, "data", "d", "", "data file (default is data.jsonl)")
	rootCmd.PersistentFlags().BoolVar(&flageh, "eh", false, "-eh for event hub")
	rootCmd.PersistentFlags().BoolVar(&flagfile, "f", false, "-f for file")

	rootCmd.Run = func(cmd *cobra.Command, args []string) {

		if flageh && !flagfile {
			if configFile != "" {
				conf = openConfigFile(configFile)
			} else {
				conf = openConfigFile("config.yaml")
			}
			data = datasource.NewDataSourceEH(
				conf.Connection,
				"test",
				conf.EventHubName)
		}
		if flagfile && !flageh {
			if dataFile != "" {
				data = datasource.NewDataSource("")
			} else {
				data = datasource.NewDataSource(dataFile)
			}
		}
		if !flagfile && !flageh {
			data = datasource.NewDataSource("")
		}
		service = model.NewService(data)
		mlog = NewLoggingService(service)
		apiSrv = api.NewApiController(mlog, data.GetKeys())

		apiSrv.Run()
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func openConfigFile(s string) Config {
	if s == "" {
		s = "config.yaml"
	}

	f, err := os.Open(s)
	if err != nil {
		processError(errors.Join(err, errors.New("open config.yaml file")))
	}
	defer f.Close()

	var config Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		processError(err)
	}
	return config

}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
