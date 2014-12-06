package settings

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Options struct {
	Dir      string `json:"dir" short:"d" long:"dir" default:"data/" description:"Directory to watch and sync"`
	ClientId string `json:"clientid" short:"c" long:"cid" description:"Client Id unique to this machine"`
	UserId   string `json:"userid" short:"u" long:"uid" description:"User Id unique to the user running using the client"`
}

func (o *Options) LoadFromJSON() {
	if _, err := os.Stat(ConfigFile); err == nil {
		if f, err := os.Open(ConfigFile); err == nil {
			config, _ := ioutil.ReadAll(f)
			options := &Options{}
			json.Unmarshal(config, &options)

			if o.ClientId == "" {
				o.ClientId = options.ClientId
			}

			if o.UserId == "" {
				o.UserId = options.UserId
			}
		}
	}
}
