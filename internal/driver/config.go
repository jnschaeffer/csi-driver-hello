package driver

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Path     string
	NodeName string `mapstructure:"node_name"`
}

func MustViperFlags(v *viper.Viper, flags *pflag.FlagSet) {
	flags.String("driver.path", "/csi/csi.sock", "CSI socket")
	if err := v.BindPFlag("driver.path", flags.Lookup("driver.path")); err != nil {
		panic(err)
	}

	flags.String("driver.node_name", "/csi/csi.sock", "CSI socket")
	if err := v.BindPFlag("driver.node_name", flags.Lookup("driver.node_name")); err != nil {
		panic(err)
	}
}
