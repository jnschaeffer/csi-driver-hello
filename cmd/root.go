package cmd

import (
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/jnschaeffer/csi-driver-hello/internal/config"
	"github.com/jnschaeffer/csi-driver-hello/internal/driver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			serve()
		},
	}

	cfgFile string
	sigCh   chan os.Signal
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /etc/csi-driver-hello/csi-driver-hello.yaml)")

	v := viper.GetViper()
	flags := rootCmd.Flags()

	driver.MustViperFlags(v, flags)

	sigCh = make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
}

func Execute() {
	rootCmd.Execute()
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("/etc/csi-driver-hello")
		viper.SetConfigType("yaml")
		viper.SetConfigName("csi-driver-hello")
	}

	// Allow populating configuration from environment
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("csi")
	viper.AutomaticEnv() // read in environment variables that match

	err := viper.ReadInConfig()

	if err == nil {
		log.Printf("using config file %s", viper.ConfigFileUsed())
	}

	err = viper.Unmarshal(&config.Config)
	if err != nil {
		log.Fatalf("unable to decode app config: %s", err)
	}
}

func serve() {
	server, err := driver.NewServer(config.Config.Driver)
	if err != nil {
		log.Fatal(err)
	}

	errCh := make(chan error, 1)

	go func() {
		if err := server.Run(); err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		log.Fatal(err)
	case <-sigCh:
		log.Print("received interrupt, stopping server")
		server.Stop()
	}
}
