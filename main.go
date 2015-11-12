package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/denisbakhtin/oauth2-example/config"
	"github.com/denisbakhtin/oauth2-example/facebook"
	"github.com/denisbakhtin/oauth2-example/google"
	"github.com/denisbakhtin/oauth2-example/linkedin"
	"github.com/denisbakhtin/oauth2-example/vk"
)

//state should be regenerated per auth request

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// generate loginURL
	fburl := facebook.GetConfig().AuthCodeURL(facebook.State)
	vkurl := vk.GetConfig().AuthCodeURL(vk.State)
	inurl := linkedin.GetConfig().AuthCodeURL(linkedin.State)
	gourl := google.GetConfig().AuthCodeURL(google.State)

	// Home page will display a button for login to Facebook
	fmt.Fprintf(w, `
		<html>
			<head>
				<title>Golang Oauth2 Login Example</title>
			</head>
			<body> 
				<h1>Golang Oauth2 Login Example</h1>
				<a href=%q><button>Login with Facebook!</button> </a><br>
				<a href=%q><button>Login with VK!</button> </a><br>
				<a href=%q><button>Login with Linkedin!</button> </a><br>
				<a href=%q><button>Login with Google!</button> </a><br>
				<br><br>
				<a href="/facebook_post"><button>Post on Facebook Page Wall!</button> </a><br>
			</body>
		</html>`, fburl, vkurl, inurl, gourl)

}

func init() {
	config.LoadConfig()
	http.HandleFunc("/", Home)
	http.HandleFunc("/facebook_callback", facebook.Callback)
	http.HandleFunc("/vk_callback", vk.Callback)
	http.HandleFunc("/linkedin_callback", linkedin.Callback)
	http.HandleFunc("/google_callback", google.Callback)

	http.HandleFunc("/facebook_post", facebook.PostOnPage)
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}
