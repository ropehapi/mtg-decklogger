package domain

import (
	"encoding/hex"
	"strings"
)

type DeckVersion struct {
	ID          string
	DeckId      string
	Version     int
	Description string
}

func NewDeckVersion(deckId string, version int, description string) (DeckVersion, error) {
	deckId = strings.TrimSpace(deckId)
	if deckId == "" {
		return DeckVersion{}, ErrDeckIdRequired
	}
	if len(deckId) != 24 {
		return DeckVersion{}, ErrDeckIdInvalid
	}
	if _, err := hex.DecodeString(deckId); err != nil {
		return DeckVersion{}, ErrDeckIdInvalid
	}
	return DeckVersion{
		DeckId:      deckId,
		Version:     version,
		Description: strings.TrimSpace(description),
	}, nil
}
