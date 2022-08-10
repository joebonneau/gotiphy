package playlistcmds

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/joebonneau/gotiphy/gotiphy/lib"
	"github.com/zmb3/spotify/v2"
)

const (
	EP              string = "EP"
	LP              string = "LP"
	minimumEPLength int    = 3
)

type UsefulAlbum struct {
	Album           *spotify.SimpleAlbum
	ActualAlbumType string
}

func SavePlaylistAlbums(playlistID string, mode string) error {
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
		albumInLibrary, err := client.UserHasAlbums(ctx, albumID)
		if err != nil {
			return err
		}
		if albumInLibrary[0] {
			continue
		}
		switch playlistItemAlbum.AlbumType {
		case "album", "compilation":
			viableOptions = append(
				viableOptions,
				UsefulAlbum{
					Album:           &playlistItemAlbum,
					ActualAlbumType: LP,
				},
			)
		case "single":
			trackInfo, err := client.GetAlbumTracks(ctx, albumID)
			if err != nil {
				return err
			}
			if len(trackInfo.Tracks) >= minimumEPLength {
				viableOptions = append(
					viableOptions,
					UsefulAlbum{
						Album:           &playlistItemAlbum,
						ActualAlbumType: EP,
					},
				)
			}
		default:
			continue
		}
	}

	if mode == "all" {
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
	tab.AppendHeader(table.Row{"#", "Album Name", "Artist(s)", "Album Type", "Release Date"})
	for i, album := range viableOptions {
		tab.AppendRow(table.Row{
			i,
			album.Album.Name,
			lib.Truncate(lib.GetArtistsString(album.Album.Artists), 40),
			album.ActualAlbumType,
			album.Album.ReleaseDate,
		})
	}
	tab.SetStyle(table.StyleColoredGreenWhiteOnBlack)
	tab.Render()

	scanner := bufio.NewScanner(os.Stdin)
	var intIndices []int
	for {
		fmt.Print("Enter the items you would like to save (separated by commas): ")
		scanner.Scan()
		response := scanner.Text()
		if len(response) != 0 {
			strIndices := strings.Split(response, ",")
			for _, idx := range strIndices {
				idx = strings.TrimSpace(idx)
				intIdx, err := strconv.Atoi(idx)
				if err != nil {
					fmt.Println("invalid selection entered")
					continue
				}
				intIndices = append(intIndices, intIdx)
			}
			break
		}
	}

	// add selected albums by index
	var selectedIDs []spotify.ID
	for _, i := range intIndices {
		selectedIDs = append(selectedIDs, viableOptions[i].Album.ID)
	}
	err = client.AddAlbumsToLibrary(ctx, selectedIDs...)
	if err != nil {
		return err
	}
	fmt.Println("Albums added to library successfully!")
	return nil
}
