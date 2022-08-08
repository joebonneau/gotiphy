package playlistcmds

import (
	"context"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/joebonneau/gotiphy/gotiphy/lib"
	"github.com/zmb3/spotify/v2"
)

const (
	EP string = "EP"
	LP string = "LP"
)

type UsefulAlbum struct {
	Album           spotify.SimpleAlbum
	ActualAlbumType string
}

func SavePlaylistAlbums(playlistID string, all bool) error {
	ctx := context.Background()
	client, err := lib.GetClient()
	if err != nil {
		return err
	}

	playlistItemPage, err := client.GetPlaylistItems(ctx, spotify.ID(playlistID))
	if err != nil {
		return err
	}

	var viableOptions []UsefulAlbum
	for _, playlistItem := range playlistItemPage.Items {
		playlistItemAlbum := playlistItem.Track.Track.Album
		albumID := playlistItemAlbum.ID
		switch playlistItemAlbum.AlbumType {
		case "album", "compilation":
			albumInLibrary, err := client.UserHasAlbums(ctx, albumID)
			if err != nil {
				return err
			}
			if !albumInLibrary[0] {
				viableOptions = append(
					viableOptions, 
					UsefulAlbum{
						Album: playlistItemAlbum, 
						ActualAlbumType: LP,
					}
				)
			}
		case "single":
			trackInfo, err := client.GetAlbumTracks(ctx, albumID)
			if err != nil {
				return err
			}
			if trackInfo.Total > 1 {
				viableOptions = append(
					viableOptions, 
					UsefulAlbum{
						Album: playlistItemAlbum, 
						ActualAlbumType: EP,
					}
				)
			}
		default:
			continue
		}
	}

	if all {
		var albumIDs []spotify.ID
		for _, album := range viableOptions {
			albumIDs = append(albumIDs, album.Album.ID)
		}
		err = client.AddAlbumsToLibrary(ctx, albumIDs...)
		if err != nil {
			return err
		}
		return nil
	}

	// If the user doesn't want to auto-add every viable option, create a table to display the options
	tab := table.NewWriter()
	tab.SetOutputMirror(os.Stdout)
	tab.AppendHeader(table.Row{"#", "Artist(s)", "Album Name", "Album Type", "Release Date"})
	for i, album := range viableOptions {
		tab.AppendRow(table.Row{i, lib.GetArtistsString(album.Album.Artists), album.ActualAlbumType, album.Album.ReleaseDate})
	}
	tab.Render()
	return nil
}
