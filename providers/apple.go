package providers

import (
	"Share_Playlist/config"
	"fmt"
	"log"
	"net/http"
)

type AppleMusic struct{}

func (am AppleMusic) CreatePlaylist(name string) (url string, err error) {
	client := config.AppleMusicClient{
		SubURI: "me/library/playlists",
		Token:  &config.Token,
		Method: http.MethodPost,
	}

	client.RequestBody = config.JsonMap{
		"name": name,
	}

	resp, err := client.MakeRequest()

	if err != nil {
		return "", err
	}

	playlistID := resp["data"].([]config.JsonMap)[0]["id"]

	log.Println(resp["data"].([]config.JsonMap))
	log.Println(resp["data"].([]config.JsonMap)[0])
	log.Println(resp["data"].([]config.JsonMap)[0]["id"])

	url = "https://music.apple.com/playlist/" + playlistID.(string)

	return
}

func (am AppleMusic) AddTracksToPlaylist(playlistID string, trackDetails []config.JsonMap) error {
	subUri := fmt.Sprintf("me/library/playlists/%s/tracks", playlistID)

	client := config.AppleMusicClient{
		SubURI: subUri,
		Token:  &config.Token,
		Method: http.MethodPost,
	}

	client.RequestBody = config.JsonMap{
		"tracks": config.JsonMap{
			"data": trackDetails,
		},
	}

	resp, err := client.MakeRequest()
	if err != nil {
		log.Println(resp)
		log.Println(err)
		return err
	}

	return nil
}

func (am AppleMusic) GetTrack(isrc string) (trackDetails config.JsonMap, err error) {

	subUri := "catalog/ng/songs?filter[isrc]=" + isrc

	log.Println(subUri)

	client := config.AppleMusicClient{
		RequestBody: nil,
		SubURI:      subUri,
		Token:       &config.Token,
		Method:      http.MethodGet,
	}

	resp, err := client.MakeRequest()

	if err != nil {
		log.Println(resp)
		log.Println(err)
		return
	}

	// TODO Iterate and compare titles
	trackDetails["id"] = resp["data"].([]config.JsonMap)[0]["id"].(string)
	trackDetails["type"] = resp["data"].([]config.JsonMap)[0]["type"].(string)

	return trackDetails, nil
}

func (am AppleMusic) GetPlaylist(playlistId string) (isrcList []string, err error) {
	subUri := fmt.Sprintf("me/library/playlists/%s/tracks?include=catalog", playlistId)

	client := config.AppleMusicClient{
		RequestBody: nil,
		SubURI:      subUri,
		Token:       nil,
		Method:      http.MethodGet,
	}
	resp, err := client.MakeRequest()

	if err != nil {
		log.Println(resp)
		log.Println(err)
		return isrcList, err
	}
	log.Println(resp)

	tracks := resp["data"].([]config.JsonMap)

	for _, track := range tracks {
		log.Println(track["relationships"].(config.JsonMap)["tracks"].(config.JsonMap)["data"].(config.JsonMap))
		trackID := track["relationships"].(config.JsonMap)["catalog"].(config.JsonMap)["data"].(config.JsonMap)["isrc"].(string)

		isrcList = append(isrcList, trackID)
	}

	return isrcList, nil

}
