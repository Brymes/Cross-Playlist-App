package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"io"
	"time"

	"io/ioutil"
	"log"
	"net/http"
)

var appleMusicBaseURI = "https://api.music.apple.com/v1/"
var Token AppleMusicToken

type AppleMusicToken struct{ Token string }

type JsonMap map[string]interface{}

type AppleMusicClient struct {
	RequestBody JsonMap
	SubURI      string
	Token       *AppleMusicToken
	Method      string
}

func (token *AppleMusicToken) ConstructToken() {

	mySigningKey := []byte("AllYourBase")

	claims := &jwt.StandardClaims{
		ExpiresAt: 150000,
		Issuer:    "test",
		IssuedAt:  time.Now().Unix(),
	}

	tkn := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	signedToken, err := tkn.SignedString(mySigningKey)

	if err != nil {
		log.Println("Error Creating Signed Apple Music Token")
		log.Fatal(err)
	}

	token.Token = signedToken
}

func (payload AppleMusicClient) MakeRequest() (result map[string]interface{}, err error) {
	// Maximum amount of retries if Unauthorized is 3
	retryCount := 0
	var req *http.Request
	var body *bytes.Buffer

	// Deferred function to recover panics
	defer func() {
		if err := recover(); err != nil {
			err = errors.New("Cannot Process Request at the Moment")
			return
		}
	}()

	uri := appleMusicBaseURI + payload.SubURI
	bearer := "Bearer " + payload.Token.Token

	if payload.RequestBody != nil {
		jsonPayload, err := json.Marshal(payload.RequestBody)

		if err != nil {
			panic(err)
		}

		err = json.NewEncoder(body).Encode(jsonPayload)

		if err != nil {
			panic(err)
		}
	}

	req, err = http.NewRequest(payload.Method, uri, body)

	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error Making Request to Apple Music")
		panic(err)
	}

	if resp.StatusCode > 401 {
		panic("")
	} else if resp.StatusCode == 401 {
		retryCount += 1
		if retryCount > 3 {
			panic("Cannot Authenticate Apple Music")
		}
		payload.Token.ConstructToken()
		result, err = payload.MakeRequest()
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error Reading Response from Apple Music")
		panic(err)
	}

	/// FIXME for empty Body response
	err = json.Unmarshal(response, &result)

	if err != nil {
		log.Println("Error Reading Response from Apple Music")
		panic(err)
	}

	return result, nil
}
