package main

import (
	"github.com/go-ini/ini"
)

func getConfig() map[string]string {
	ret := make(map[string]string)
	if !Exists(ConfigFile) {
		return ret
	}
	cfg, err := ini.Load(ConfigFile)
	if err != nil {
		Logger.Error(err.Error())
		return ret

	}

	for _, key := range ConfigKey {
		ret[key] = cfg.Section("").Key(key).String()
	}

	return ret
}

var Config = getConfig()
