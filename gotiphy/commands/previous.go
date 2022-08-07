package commands

import (
	"context"

	"github.com/joebonneau/gotiphy/gotiphy/lib"
)

func PreviousTrack() error {
	ctx := context.Background()

	client, err := lib.GetClient()
	if err != nil {
		return err
	}

	err = client.Previous(ctx)
	if err != nil {
		return err
	}

	err = lib.GetAndDisplayCurrentPlayback(ctx, *client)
	if err != nil {
		return err
	}

	return nil
}
