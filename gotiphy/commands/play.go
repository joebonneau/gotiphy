package commands

import (
	"context"

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

	err = lib.GetAndDisplayCurrentPlayback(ctx, *client)
	if err != nil {
		return err
	}

	return nil
}
