package main

import "Share_Playlist/config"

func init() {
	config.InitDb()
	config.InitSpotifyClient()
	config.Token.ConstructToken()
}

func main() {
	Server()
}
