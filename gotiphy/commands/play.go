package commands

import (
	"context"
	"fmt"

	"github.com/joebonneau/gotiphy/gotiphy/lib"
	spotify "github.com/zmb3/spotify/v2"
)

func StartPlayback(uri string) error {
	ctx := context.Background()

	client, err := lib.GetClient()
	if err != nil {
		return err
	}

	if uri != "" {
		uris := []spotify.URI{spotify.URI(uri)}
		options := spotify.PlayOptions{
			URIs: uris,
		}
		err = client.PlayOpt(ctx, &options)
		if err != nil {
			return err
		}
	} else {
		err = client.Play(ctx)
		if err != nil {
			return err
		}
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
