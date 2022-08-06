package lib

import (
	"context"
	"fmt"
	"os"
	"time"

	spotify "github.com/zmb3/spotify/v2"
	spauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const RedirectURI = "http://localhost:8080/callback"

var (
	scopes = []string{
		spauth.ScopeUserModifyPlaybackState,
		spauth.ScopeUserReadPlaybackState,
		spauth.ScopeUserLibraryRead,
		spauth.ScopeUserLibraryModify,
		spauth.ScopePlaylistReadPrivate,
		spauth.ScopePlaylistReadCollaborative,
		spauth.ScopePlaylistModifyPublic,
		spauth.ScopePlaylistModifyPrivate,
	}
	Auth = spauth.New(spauth.WithRedirectURL(RedirectURI), spauth.WithScopes(scopes...))
)

func GetClient() (*spotify.Client, error) {
	ctx := context.TODO()

	t, err := os.ReadFile("refresh.token")
	if err != nil {
		return nil, fmt.Errorf("file 'refresh.token' does not exist. Execute the 'gotiphy auth' command to generate")
	}

	token := new(oauth2.Token)
	token.Expiry = time.Now().Add(time.Second * -5)
	token.RefreshToken = string(t)

	// use the token to get an authenticated client
	client := spotify.New(Auth.Client(ctx, token))

	// set the new refresh token for the next request
	newToken, err := client.Token()
	if err != nil {
		return nil, err
	}

	err = os.WriteFile("refresh.token", []byte(newToken.RefreshToken), 0644)
	if err != nil {
		return nil, err
	}

	return client, nil
}
