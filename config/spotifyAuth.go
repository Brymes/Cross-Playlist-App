package config

import (
	"context"
	"github.com/zmb3/spotify/v2"
	"github.com/zmb3/spotify/v2/auth"
	"log"
	"os"

	"golang.org/x/oauth2/clientcredentials"
)

var SpotifyClient *spotify.Client
var UserID = os.Getenv("SPOTIFY_USER_ID") //TODO ClientID?

func InitSpotifyClient() {

	ctx := context.Background()

	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := config.Token(ctx)

	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)

	SpotifyClient = spotify.New(httpClient)
}
