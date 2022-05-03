package providers

import (
	"Share_Playlist/config"
	"context"
	"github.com/zmb3/spotify/v2"
)

type Spotify struct{}

func (s Spotify) GetPlaylistSpotify(playlistId string) (isrcList []string, err error) {

	id := spotify.ID(playlistId)
	//tracks, err := spotify.Client.GetPlaylistTracks(context.Background(), context.Background(), id)
	tracks, err := config.SpotifyClient.GetPlaylistTracks(context.Background(), id)

	//isrcList = make([]string, tracks.Total)

	if err != nil {
		return nil, err
	}

	for _, track := range tracks.Tracks {
		isrcList = append(isrcList, track.Track.ExternalIDs["isrc"])
	}

	return isrcList, nil
}

func (s Spotify) GetTrack(isrc string) (trackID spotify.ID, err error) {

	res, err := config.SpotifyClient.Search(context.Background(), isrc, spotify.SearchTypeTrack)

	trackID = res.Tracks.Tracks[0].ID

	return trackID, nil
}

func (s Spotify) CreatePlaylist(name string) (id spotify.ID, err error) {
	playlist, err := config.SpotifyClient.CreatePlaylistForUser(context.Background(), config.UserID, name, "", true, false)

	return playlist.ID, nil
}

func (s Spotify) AddTracksToPlaylist(playlistID spotify.ID, trackIDs ...spotify.ID) (link string, err error) {
	link, err = config.SpotifyClient.AddTracksToPlaylist(context.Background(), playlistID, trackIDs...)

	if err != nil {
		return "", err
	}

	return link, nil
}
