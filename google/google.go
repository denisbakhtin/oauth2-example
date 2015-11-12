package google

import (
	"fmt"
	"net/http"

	"github.com/denisbakhtin/oauth2-example/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	goauth2 "google.golang.org/api/oauth2/v2"
)

//state should be regenerated per auth request
var (
	State = "google_random_csrf_string"
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
	service, err := goauth2.New(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	uService := goauth2.NewUserinfoService(service)
	gouser, err := uService.Get().Do()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Username: %s %s<br>ID: %s<br>Email: %s<br>Picture: %s<br>", gouser.GivenName, gouser.FamilyName, gouser.Id, gouser.Email, gouser.Picture)
}

func GetConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.GetConfig().Google.ClientID,     // change this to yours
		ClientSecret: config.GetConfig().Google.ClientSecret, //change this to yours
		RedirectURL:  config.GetConfig().Google.RedirectURL,  // change this to your webserver adddress
		Scopes:       []string{goauth2.PlusLoginScope, goauth2.PlusMeScope, goauth2.UserinfoEmailScope, goauth2.UserinfoProfileScope},
		Endpoint:     google.Endpoint,
	}
}
