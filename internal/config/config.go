package config

import (
	"github.com/jnschaeffer/csi-driver-hello/internal/driver"
	"github.com/jnschaeffer/csi-driver-hello/internal/manager"
)

var Config struct {
	Driver  driver.Config
	Manager manager.Config
}
