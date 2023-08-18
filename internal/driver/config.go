package driver

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Path string
}

func MustViperFlags(v *viper.Viper, flags *pflag.FlagSet) {
	flags.String("driver.path", "/csi/csi.sock", "CSI socket")
	if err := v.BindPFlag("driver.path", flags.Lookup("driver.path")); err != nil {
		panic(err)
	}
}
