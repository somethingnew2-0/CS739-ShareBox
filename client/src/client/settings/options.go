package settings

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"os"
)

type Options struct {
	Dir       string `json:"dir" short:"d" long:"dir" default:"data/" description:"Directory to watch and sync"`
	Username  string `json:"username" short:"u" long:"username" default:"test" description:"Username for user login, unique to the userid"`
	Password  string `json:"password" short:"p" long:"password" default:"test" description:"Password to encrpyt file with"`
	Hash      []byte `json:"hash" short:"h" long:"hash" description:"Calcuated password hash to encrpyt file with"`
	ClientId  string `json:"clientid" short:"c" long:"cid" description:"Client Id unique to this machine"`
	UserId    string `json:"userid" long:"uid" description:"User Id unique to the user running using the client"`
	AuthToken string `json:"authToken" short:"t" long:"token" description:"The session token required to change anything on the server"`
}

func (o *Options) Load() {
	if _, err := os.Stat(ConfigFile); err == nil {
		if f, err := os.Open(ConfigFile); err == nil {
			config, _ := ioutil.ReadAll(f)
			options := &Options{}
			if err := json.Unmarshal(config, &options); err == nil {
				if o.Dir == "" {
					o.Dir = options.Dir
				}
				if o.Username == "" {
					o.Username = options.Username
				}
				if len(o.Hash) == 0 {
					o.Hash = options.Hash
				}
				if o.ClientId == "" {
					o.ClientId = options.ClientId
				}
				if o.UserId == "" {
					o.UserId = options.UserId
				}
				if o.AuthToken == "" {
					o.AuthToken = options.AuthToken
				}
			}
		}
	}
}

func (o *Options) HashPassword() {
	if o.Password != "" && len(o.Hash) == 0 {
		o.Hash, _ = bcrypt.GenerateFromPassword([]byte(o.Password), 10)
		o.Password = ""
	}
}

func (o *Options) Save() {
	o.HashPassword()

	config, _ := json.Marshal(o)
	ioutil.WriteFile(ConfigFile, config, 0666)
}
