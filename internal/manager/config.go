package manager

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Path string
}

func MustViperFlags(v *viper.Viper, flags *pflag.FlagSet) {
	flags.String("manager.path", "/tmp/manager", "Manager path")
	if err := v.BindPFlag("manager.path", flags.Lookup("manager.path")); err != nil {
		panic(err)
	}
}
