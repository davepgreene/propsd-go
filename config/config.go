package config

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

var index = map[string]interface{}{
	"path":		"index.json",
	"interval":  	30000,
	"region": 	"us-east-1",
}

var service = map[string]interface{} {
	"host": 	"127.0.0.1",
	"port":		9100,
}

var log = map[string]interface{}{
	"level": 	logrus.InfoLevel,
	"json":     	true,
	"requests": 	true,
}

var consul = map[string]interface{}{
	"host": 	"127.0.0.1",
	"port": 	8500,
	"secure": 	false,
}

var tokend = map[string]interface{}{
	"host":		"127.0.0.1",
	"port":		4500,
	"interval":	300000,
}

var metadata = map[string]interface{}{
	"host": 	"169.254.169.254",
	"interval": 	30000,
	"timeout":	500,
	"version":	"latest",
}

var tags = map[string]interface{}{
	"interval":	300000,
}

// Defaults generates a set of default configuration options
func Defaults() {
	viper.SetDefault("index", index)
	viper.SetDefault("service", service)
	viper.SetDefault("log", log)
	viper.SetDefault("consul", consul)
	viper.SetDefault("tokend", tokend)
	viper.SetDefault("metadata", metadata)
	viper.SetDefault("tags", tags)
}
