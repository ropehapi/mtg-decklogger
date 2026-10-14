package mongodb

import (
	"context"

	"github.com/ropehapi/mtg-decklogger/internal/deck/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	collection *mongo.Collection
}

type deckDocument struct {
	Name        string `bson:"name"`
	Format      string `bson:"format"`
	Description string `bson:"description"`
}

func NewRepository(collection *mongo.Collection) *Repository {
	return &Repository{collection: collection}
}

func (repository *Repository) Create(ctx context.Context, deck domain.Deck) error {
	document := deckDocument{
		Name:        deck.Name,
		Format:      deck.Format,
		Description: deck.Description,
	}

	_, err := repository.collection.InsertOne(ctx, document)
	return err
}

func (r *Repository) Find(ctx context.Context) ([]domain.Deck, error) {
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documents []deckDocument
	if err := cursor.All(ctx, &documents); err != nil {
		return nil, err
	}

	decks := make([]domain.Deck, 0, len(documents))
	for _, document := range documents {
		decks = append(decks, domain.Deck{
			Name:        document.Name,
			Format:      document.Format,
			Description: document.Description,
		})
	}

	return decks, nil
}
