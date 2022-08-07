package lib

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify/v2"
)

func GetAndDisplayCurrentPlayback(ctx context.Context, client spotify.Client) error {
	nowPlaying, err := client.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return err
	}
	item := nowPlaying.Item
	artistsString := GetArtistsString(item.Artists)
	fmt.Printf("Now playing: %s by %s from the album %s (%v)", item.Name, artistsString, item.Album.Name, item.Album.ReleaseDateTime().Year())
	return nil
}
