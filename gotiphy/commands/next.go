package commands

import (
	"context"
	"fmt"

	"github.com/joebonneau/gotiphy/gotiphy/lib"
)

func NextTrack() error {
	ctx := context.Background()

	client, err := lib.GetClient()
	if err != nil {
		return err
	}

	err = client.Next(ctx)
	if err != nil {
		return err
	}

	nowPlaying, err := client.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return err
	}
	item := nowPlaying.Item
	artistsString := lib.GetArtistsString(item.Artists)
	fmt.Printf("Now playing: %s by %s from the album %s (%v)", item.Name, artistsString, item.Album.Name, item.Album.ReleaseDateTime().Year())
	return nil
}
