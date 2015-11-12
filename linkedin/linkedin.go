package linkedin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/denisbakhtin/oauth2-example/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
)

//state should be regenerated per auth request
var (
	State = "linkedin_random_csrf_string"
)

func Callback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//check state validity, see url := Config.AuthCodeURL(state) above
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
	req, err := http.NewRequest("GET", "https://api.linkedin.com/v1/people/~:(email-address,first-name,last-name,id,headline)?format=json", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Header.Set("Bearer", token.AccessToken)
	response, err := client.Do(req)

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

	var inuser struct {
		Id        string
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Headline  string
		Email     string `json:"emailAddress"`
	}

	err = json.Unmarshal(str, &inuser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Username: %s %s<br>ID: %s<br>Email: %s<br>Headline: %s<br>", inuser.FirstName, inuser.LastName, inuser.Id, inuser.Email, inuser.Headline)
}

func GetConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.GetConfig().Linkedin.ClientID,     // change this to yours
		ClientSecret: config.GetConfig().Linkedin.ClientSecret, //change this to yours
		RedirectURL:  config.GetConfig().Linkedin.RedirectURL,  // change this to your webserver adddress
		Scopes:       []string{"r_basicprofile", "r_emailaddress"},
		Endpoint:     linkedin.Endpoint,
	}
}
