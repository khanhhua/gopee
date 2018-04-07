package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/khanhhua/gopee/dao"
)

// Authorize Dropbox OAuth2 callback
func Authorize(w http.ResponseWriter, r *http.Request) {
	var code string

	if code = r.URL.Query().Get("code"); len(code) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("Auth code missing"))

		return
	}

	form := url.Values{}
	form.Add("code", code)
	form.Add("grant_type", "authorization_code")
	form.Add("client_id", "j4365xi2ynl3zri")
	form.Add("client_secret", "7e9j352ahi7hu8v")
	form.Add("redirect_uri", "http://localhost:8888/auth/dropbox")

	type AuthResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		AccountID   string `json:"account_id"`
		UID         string `json:"uid"`
	}

	if resp, err := http.PostForm("https://api.dropboxapi.com/oauth2/token", form); err == nil {
		var data []byte
		if data, err = ioutil.ReadAll(resp.Body); err != nil {
			http.Error(w, "Server error", 500)

			return
		}

		fmt.Printf("Response: %s", data)
		response := AuthResponse{}
		if json.Unmarshal(data, &response) != nil {
			http.Error(w, "Server error", 500)

			return
		}

		client := dao.Client{
			ClientKey:          response.UID,
			ClientDomain:       "",
			DropboxAccountID:   response.AccountID,
			DropboxAccessToken: response.AccessToken}
		if client, err = dao.CreateClient(client); err != nil {
			fmt.Printf("Error: %v", err)
			http.Error(w, "Database error", 500)
			return
		}

		fmt.Printf("Generated client: %v\n", client)
		http.Redirect(w, r, "http://localhost:8888/console", 301)
	} else {
		http.Error(w, "Authorization error", 403)
	}
}
