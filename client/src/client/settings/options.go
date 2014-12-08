package settings

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"os"
)

type Options struct {
	Dir      string `json:"dir" short:"d" long:"dir" default:"data/" description:"Directory to watch and sync"`
	Password string `json:"password" short:"p" long:"password" default:"test" description:"Password to encrpyt file with"`
	Hash     []byte `json:"hash" short:"h" long:"hash" description:"Calcuated password hash to encrpyt file with"`
	ClientId string `json:"clientid" short:"c" long:"cid" description:"Client Id unique to this machine"`
	UserId   string `json:"userid" short:"u" long:"uid" description:"User Id unique to the user running using the client"`
}

func (o *Options) Load() {
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

func (o *Options) HashPassword() {
	o.Hash, _ = bcrypt.GenerateFromPassword([]byte(o.Password), 10)
	o.Password = ""
}

func (o *Options) Save() {
	o.HashPassword()

	config, _ := json.Marshal(o)
	ioutil.WriteFile(ConfigFile, config, 0666)
}
