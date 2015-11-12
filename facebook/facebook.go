package facebook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/denisbakhtin/oauth2-example/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

//state should be regenerated per auth request
var (
	State = "facebook_random_csrf_string"
)

func Callback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
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

	client := GetConfig().Client(oauth2.NoContext, token)
	response, err := client.Get(fmt.Sprintf("https://graph.facebook.com/me?access_token=%s&fields=name,email,birthday", token.AccessToken))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer response.Body.Close()
	str, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var fauser struct {
		Id       string
		Name     string
		Email    string
		Birthday string
	}

	err = json.Unmarshal([]byte(str), &fauser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Username: %s<br>ID: %s<br>Birthday: %s<br>Email: %s<br>", fauser.Name, fauser.Id, fauser.Birthday, fauser.Email)

	img := fmt.Sprintf("https://graph.facebook.com/%s/picture?width=180&height=180", fauser.Id)

	fmt.Fprintf(w, "Photo is located at %s<br>", img)
	fmt.Fprintf(w, "<img src=%q>", img)
}

func PostOnPage(w http.ResponseWriter, r *http.Request) {
	//see http://stackoverflow.com/questions/17197970/facebook-permanent-page-access-token
	//for info on obtaining upexpirable page access token
	//also https://developers.facebook.com/docs/graph-api/reference/v2.5/page/feed for api description

	token := &oauth2.Token{
		AccessToken: config.GetConfig().Facebook.Token, //page access token
	}
	client := GetConfig().Client(oauth2.NoContext, token)
	response, err := client.Post(
		fmt.Sprintf(
			"https://graph.facebook.com/v2.5/dengraphapitest/feed?access_token=%s&link=%s&name=%s&caption=%s&description=%s&message=%s",
			token.AccessToken,
			url.QueryEscape("http://google.com"),
			url.QueryEscape("Link Name"),
			url.QueryEscape("Link Caption"),
			url.QueryEscape("Link Description"),
			url.QueryEscape("Post message"),
			//add picture field for image link
		),
		"application/json",
		nil,
	)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Fprintf(w, "Successfully posted, response: %+v\n", response)
	body, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	fmt.Fprintf(w, "Body: %s\n", body)
}

func GetConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.GetConfig().Facebook.ClientID,     // change this to yours
		ClientSecret: config.GetConfig().Facebook.ClientSecret, //change this to yours
		RedirectURL:  config.GetConfig().Facebook.RedirectURL,  // change this to your webserver adddress
		Scopes:       []string{"email", "user_birthday", "user_location", "user_about_me"},
		Endpoint:     facebook.Endpoint,
	}
}
