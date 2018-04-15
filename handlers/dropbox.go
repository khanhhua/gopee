package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/khanhhua/gopee/dao"
)

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	AccountID   string `json:"account_id"`
	UID         string `json:"uid"`
	Error       string `json:"error"`
}

// Authorize Dropbox OAuth2 callback
func Authorize(w http.ResponseWriter, r *http.Request) {
	BASE_URL := os.Getenv("BASE_URL")
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
	form.Add("redirect_uri", fmt.Sprintf("%s/auth/dropbox", BASE_URL))

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
		http.Redirect(w, r, fmt.Sprintf("%s/console", BASE_URL), 301)
	} else {
		http.Error(w, "Authorization error", 403)
	}
}

func GetToken(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetToken...")

	BASE_URL := os.Getenv("BASE_URL")
	var payload struct {
		Code string `json:"code"`
	}

	if data, err := ioutil.ReadAll(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	} else if err := json.Unmarshal(data, &payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	fmt.Printf("Authorization code: %s\n", payload.Code)
	if len(payload.Code) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Required 'code' is missing"))
		return
	}

	redirectURI := fmt.Sprintf("%s/console", BASE_URL)
	if response, err := validateCode(payload.Code, redirectURI); err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
	} else {
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

		// TODO JWT tokenize instead of the UID plainly
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.UID)
	}
}

func validateCode(code string, redirectURI string) (response AuthResponse, err error) {
	ClientID := os.Getenv("DROPBOX_CLIENT_ID")
	ClientSecret := os.Getenv("DROPBOX_CLIENT_SECRET")

	form := url.Values{}
	form.Add("code", code)
	form.Add("grant_type", "authorization_code")
	form.Add("client_id", ClientID)
	form.Add("client_secret", ClientSecret)
	form.Add("redirect_uri", redirectURI)

	if resp, httpErr := http.PostForm("https://api.dropboxapi.com/oauth2/token", form); httpErr != nil {
		err = httpErr
		return
	} else {
		var data []byte
		if data, err = ioutil.ReadAll(resp.Body); err != nil {
			return
		}

		fmt.Printf("Response: %s", data)
		response = AuthResponse{}
		if err = json.Unmarshal(data, &response); err != nil {
			return
		} else if len(response.Error) != 0 {
			err = errors.New(response.Error)
			return
		}

		return
	}
}
