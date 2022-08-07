package rdmcommands

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/joebonneau/gotiphy/gotiphy/lib"
	"github.com/zmb3/spotify/v2"
)

func RandomUserAlbum(action string) error {
	ctx := context.Background()
	client, err := lib.GetClient()
	if err != nil {
		return err
	}

	var userAlbums []spotify.SavedAlbum
	offset := 0
	for {
		savedAlbums, err := client.CurrentUsersAlbums(ctx, spotify.Limit(50), spotify.Offset(offset))
		if err != nil {
			return err
		}
		userAlbums = append(userAlbums, savedAlbums.Albums...)
		if len(savedAlbums.Albums) < 50 {
			// since the limit is set to 50, we know that if the number of returned albums
			// is less than 50 that we're done retrieving data
			break
		} else {
			offset += 50
		}
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(userAlbums))
	selectedAlbum := userAlbums[randomIndex]

	albumURI := spotify.URI(selectedAlbum.URI)
	switch action {
	case "play":
		options := &spotify.PlayOptions{
			PlaybackContext: &albumURI,
		}

		err = client.PlayOpt(ctx, options)
		if err != nil {
			return err
		}

		// a short delay is required here to prevent an error
		time.Sleep(time.Millisecond * 700)
		err = lib.GetAndDisplayCurrentPlayback(ctx, *client)
		if err != nil {
			return err
		}
		return nil
	case "queue":
		var albumSongsIDs []spotify.ID
		for _, track := range selectedAlbum.FullAlbum.Tracks.Tracks {
			albumSongsIDs = append(albumSongsIDs, track.ID)
		}
		for _, trackID := range albumSongsIDs {
			err = client.QueueSong(ctx, trackID)
			if err != nil {
				return err
			}
		}
		albumArtists := lib.GetArtistsString(selectedAlbum.Artists)
		fmt.Printf("Added %v songs to the queue from the album %s by %s", len(albumSongsIDs), selectedAlbum.Name, albumArtists)
		return nil
	}
	return nil
}
