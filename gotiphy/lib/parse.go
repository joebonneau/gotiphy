package lib

import (
	"strings"

	"github.com/zmb3/spotify/v2"
)

func GetArtistsString(artists []spotify.SimpleArtist) string {
	var result []string
	for _, a := range artists {
		result = append(result, a.Name)
	}
	return strings.Join(result, ", ")
}

func Truncate(str string, numChars int) string {
	if numChars != 0 && len(str) > numChars {
		return str[:numChars] + "..."
	}
	return str
}
