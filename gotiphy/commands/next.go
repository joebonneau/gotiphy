package commands

import (
	"context"

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

	err = lib.GetAndDisplayCurrentPlayback(ctx, *client)
	if err != nil {
		return err
	}

	return nil
}
