package commands

import (
	"context"
	"fmt"

	"github.com/joebonneau/gotiphy/gotiphy/lib"
)

func NowPlaying() error {
	ctx := context.Background()

	client, err := lib.GetClient()
	if err != nil {
		return err
	}

	nowPlaying, err := client.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return err
	}

	if !nowPlaying.Playing {
		fmt.Println("Nothing is currently playing.")
		return nil
	}

	item := nowPlaying.Item
	artistsString := lib.GetArtistsString(item.Artists)
	fmt.Printf("Now playing: %s by %s from the album %s (%v)", item.Name, artistsString, item.Album.Name, item.Album.ReleaseDateTime().Year())
	return nil
}
