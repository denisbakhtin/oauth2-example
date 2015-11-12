package config

import (
	"encoding/json"
	"io/ioutil"
)

//OauthConfig contains oauth login info
type Config struct {
	Facebook OauthApp
	Google   OauthApp
	Linkedin OauthApp
	Vk       OauthApp
}

//OauthApp contains oauth application data
type OauthApp struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
	Token        string `json:"token"`
}

var (
	config *Config
)

func LoadConfig() {
	data, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		panic(err)
	}
	config = &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}
}

func GetConfig() *Config {
	return config
}
