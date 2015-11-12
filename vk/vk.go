package vk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/denisbakhtin/oauth2-example/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"
)

const (
	State = "vkrandomstate"
)

func Callback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//check state validity
	state_check := r.FormValue("state")
	if State != state_check {
		http.Error(w, fmt.Sprintf("Wrong state string: Expected %s, got %s. Please, try again", State, state_check), http.StatusBadRequest)
		return
	}

	token, err := GetConfig().Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := token.Extra("email")
	userId := int64(token.Extra("user_id").(float64))
	//if you need to invoke vk api, create a client, make your requests
	client := GetConfig().Client(oauth2.NoContext, token)
	response, err := client.Get(fmt.Sprintf("https://api.vk.com/method/users.get?access_token=%s&user_id=%d", token.AccessToken, userId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer response.Body.Close()
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type Response struct {
		Uid       int64
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	var vkuser struct {
		Response []Response
	}
	if err := json.Unmarshal(buf, &vkuser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(vkuser.Response) > 0 {
		fmt.Fprintf(w, "User Name: %s %s<br> ID: %d<br>Email: %s<br>", vkuser.Response[0].FirstName, vkuser.Response[0].LastName, vkuser.Response[0].Uid, email)
	} else {
		fmt.Fprint(w, "Данные пользователя не получены")
	}
}

func GetConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.GetConfig().Vk.ClientID,     // change this to yours
		ClientSecret: config.GetConfig().Vk.ClientSecret, //change this to yours
		RedirectURL:  config.GetConfig().Vk.RedirectURL,  // change this to your webserver adddress
		Scopes:       []string{"email"},
		Endpoint:     vk.Endpoint,
	}
}
