package domain

import "strings"

type DeckVersion struct {
	ID          string
	DeckId      string
	Version     int
	Description string
}

func NewDeckVersion(deckId string, version int, description string) (DeckVersion, error){
	deckId = strings.TrimSpace(deckId)
	if deckId == "" {
		return DeckVersion{}, ErrDeckIdRequired
	}
	return DeckVersion{
		DeckId: deckId,
		Version: version,
		Description: description,
	}, nil
}