package main

import (
	"golang.org/x/oauth2/clientcredentials"
	"log"
	"context"
	"net/url"
	"net/http"
	"fmt"
	"bytes"
	"io"
	"os"
	"strings"
	"encoding/json"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loadin2g .env file")
	}

	domain := os.Getenv("AUTH0_DOMAIN")
	
	cfg := &clientcredentials.Config{
		ClientID:       os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret:   os.Getenv("AUTH0_CLIENT_SECRET"),
		Scopes:         []string{"read:users", "read:user_idp_tokens", "create:users"},
		TokenURL:       fmt.Sprintf("https://%s/oauth/token", domain),
		EndpointParams: url.Values{"audience": {os.Getenv("AUTH0_AUDIENCE")}},
	}

	ctx := context.Background()
	t, e := cfg.Token(ctx)
	log.Println(e)
	log.Printf("%+v",t)


	body, err := request("GET", fmt.Sprintf("https://%s/api/v2/users",domain), nil, t.AccessToken )
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s", body)


	// Create User
	values := struct{
		Email string `json:"email"`
		Password string `json:"password"`
		Connection string `json:"connection"`
	}{
		Email: "test@example.com",
		Password: "hogehoge",
		Connection: "Username-Password-Authentication",
	}

	
	b, _ := json.Marshal(values)

	log.Println(string(b))

	body, err = request("POST", fmt.Sprintf("https://%s/api/v2/users", domain), strings.NewReader(string(b)), t.AccessToken)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s", body)
}


func request(method, url string, values io.Reader, token string) ([]byte, error) {
	req, err := http.NewRequest(
		method,
		url,
		values,
	)
	if err != nil {
		return nil, err
	}

	if method == "POST" {
		req.Header.Set("Content-Type", "application/json")

	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}	

	// Read the response body
	buf := new(bytes.Buffer)
	io.Copy(buf, res.Body)
	res.Body.Close()
	return buf.Bytes(), nil
}
