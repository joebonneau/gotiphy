package commands

import (
	"context"

	"github.com/joebonneau/gotiphy/gotiphy/lib"
)

func PausePlayback() error {
	ctx := context.Background()

	client, err := lib.GetClient()
	if err != nil {
		return err
	}

	err = client.Pause(ctx)
	if err != nil {
		return err
	}

	return nil
}
