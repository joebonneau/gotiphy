package main

import (
	"log"
	"os"

	"github.com/joebonneau/gotiphy/gotiphy/commands"
	rdmcommands "github.com/joebonneau/gotiphy/gotiphy/commands/random"
	"github.com/urfave/cli/v2"
)

var rdmFlags []cli.Flag = []cli.Flag{
	&cli.BoolFlag{
		Name:    "play",
		Aliases: []string{"p"},
		Usage:   "Play the randomly selected item",
		Value:   true,
	},
	&cli.BoolFlag{
		Name:    "queue",
		Aliases: []string{"q"},
		Usage:   "Queue the randomly selected item",
	},
}

func main() {
	app := &cli.App{
		Name:    "gotiphy",
		Usage:   "A CLI written in Go to control Spotify player with additional playlist creation functionality",
		Version: "v0.1",
		Commands: []*cli.Command{
			{
				Name:  "auth",
				Usage: "Authenticates the user and generates auth token",
				Action: func(cCtc *cli.Context) error {
					err := commands.Authenticate()
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:  "now",
				Usage: "Display information about current playback",
				Action: func(cCtx *cli.Context) error {
					err := commands.NowPlaying()
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name:  "next",
				Usage: "Skip playback to next track in user queue",
				Action: func(cCtx *cli.Context) error {
					err := commands.NextTrack()
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name:  "prev",
				Usage: "Skip playback to previous track in user queue",
				Action: func(cCtx *cli.Context) error {
					err := commands.PreviousTrack()
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name:  "play",
				Usage: "Play a specific item or resume current playback",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "uri",
						Aliases: []string{"u"},
						Usage:   "Play a specific Spotify URI",
					},
				},
				Action: func(cCtx *cli.Context) error {
					uri := cCtx.String("uri")
					err := commands.StartPlayback(uri)
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name:  "pause",
				Usage: "Pauses current playback",
				Action: func(cCtx *cli.Context) error {
					err := commands.PausePlayback()
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name:    "random",
				Aliases: []string{"rdm"},
				Usage:   "Options for playing or queueing random media",
				Subcommands: []*cli.Command{
					{
						Name:    "album",
						Aliases: []string{"a"},
						Usage:   "Selects an album at random from user library",
						Flags:   rdmFlags,
						Action: func(cCtx *cli.Context) error {
							if cCtx.NArg() > 1 {
								log.Fatal("too many flags were specified")
							}
							action := "play"
							if cCtx.Bool("queue") {
								action = "queue"
							}
							err := rdmcommands.RandomUserAlbum(action)
							if err != nil {
								log.Fatal(err)
							}
							return nil
						},
					},
					{
						Name:    "playlist",
						Aliases: []string{"p", "pl"},
						Usage:   "Selects a playlist at random from user library",
						Flags:   rdmFlags,
						Action: func(cCtx *cli.Context) error {
							if cCtx.NArg() > 1 {
								log.Fatal("too many flags were specified")
							}
							action := "play"
							if cCtx.Bool("queue") {
								action = "queue"
							}
							err := rdmcommands.RandomUserPlaylist(action)
							if err != nil {
								log.Fatal(err)
							}
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
