package rdmcommands

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/joebonneau/gotiphy/gotiphy/lib"
	"github.com/zmb3/spotify/v2"
)

func RandomUserPlaylist(action string) error {
	ctx := context.Background()
	client, err := lib.GetClient()
	if err != nil {
		return err
	}

	var userPlaylists []spotify.SimplePlaylist
	offset := 0
	for {
		savedPlaylists, err := client.CurrentUsersPlaylists(ctx, spotify.Limit(50), spotify.Offset(offset))
		if err != nil {
			return err
		}
		userPlaylists = append(userPlaylists, savedPlaylists.Playlists...)
		if len(savedPlaylists.Playlists) < 50 {
			// since the limit is set to 50, we know that if the number of returned
			// playlists is less than 50 that we're done retrieving data
			break
		} else {
			offset += 50
		}
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(userPlaylists))
	selectedPlaylist := userPlaylists[randomIndex]
	playlistURI := spotify.URI(selectedPlaylist.URI)
	switch action {
	case "play":
		options := &spotify.PlayOptions{
			PlaybackContext: &playlistURI,
		}

		err = client.PlayOpt(ctx, options)
		if err != nil {
			return err
		}

		// a short delay is required here to prevent an error
		time.Sleep(time.Millisecond * 700)
		fmt.Printf("Selected playlist: %s\n", selectedPlaylist.Name)
		err = lib.GetAndDisplayCurrentPlayback(ctx, *client)
		if err != nil {
			return err
		}
		return nil
	case "queue":
		playlistItemPages, err := client.GetPlaylistItems(ctx, selectedPlaylist.ID)
		if err != nil {
			return err
		}
		var playlistItemIDs []spotify.ID
		for _, playlistItem := range playlistItemPages.Items {
			playlistItemTrack := playlistItem.Track
			var playlistItemID spotify.ID
			if playlistItemTrack.Track != nil {
				playlistItemID = playlistItemTrack.Track.ID
			} else if playlistItemTrack.Episode != nil {
				playlistItemID = playlistItemTrack.Episode.ID
			} else {
				continue
			}
			playlistItemIDs = append(playlistItemIDs, playlistItemID)
		}
		for _, itemID := range playlistItemIDs {
			err = client.QueueSong(ctx, itemID)
			if err != nil {
				return err
			}
		}
		fmt.Printf("Added %v songs to the queue from the playlist %s", len(playlistItemIDs), selectedPlaylist.Name)
		return nil
	}
	return nil

}
