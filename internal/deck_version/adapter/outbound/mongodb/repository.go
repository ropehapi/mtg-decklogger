package mongodb

import (
	"context"

	"github.com/ropehapi/mtg-decklogger/internal/deck_version/domain.go"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository struct {
	collection *mongo.Collection
}

type deckVersionDocument struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	DeckID      bson.ObjectID `bson:"deck_id"`
	Version     int           `bson:"version"`
	Description string        `bson:"description"`
}

func NewRepository(collection *mongo.Collection) *Repository {
	return &Repository{collection: collection}
}

func (repository *Repository) EnsureIndexes(ctx context.Context) error {
	_, err := repository.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "deck_id", Value: 1},
			{Key: "version", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	return err
}

func (repository *Repository) Create(ctx context.Context, deck domain.DeckVersion) error {
	deckID, err := bson.ObjectIDFromHex(deck.DeckId)
	if err != nil {
		return err
	}

	document := deckVersionDocument{
		DeckID:      deckID,
		Version:     deck.Version,
		Description: deck.Description,
	}

	_, err = repository.collection.InsertOne(ctx, document)
	if mongo.IsDuplicateKeyError(err) {
		return domain.ErrDeckVersionAlreadyExists
	}

	return err
}

// func (r *Repository) Find(ctx context.Context) ([]domain.Deck, error) {
// 	cursor, err := r.collection.Find(ctx, bson.D{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	var documents []deckVersionDocument
// 	if err := cursor.All(ctx, &documents); err != nil {
// 		return nil, err
// 	}

// 	decks := make([]domain.Deck, 0, len(documents))
// 	for _, document := range documents {
// 		decks = append(decks, domain.Deck{
// 			ID:          document.ID.Hex(),
// 			Name:        document.Name,
// 			Format:      document.Format,
// 			Description: document.Description,
// 		})
// 	}

// 	return decks, nil
// }

// func (r *Repository) FindById(ctx context.Context, id string) (domain.Deck, error) {
// 	objectID, err := bson.ObjectIDFromHex(id)
// 	if err != nil {
// 		return domain.Deck{}, err
// 	}

// 	var document deckVersionDocument
// 	if err := r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&document); err != nil {
// 		return domain.Deck{}, err
// 	}

// 	return domain.Deck{
// 		Name:        document.Name,
// 		Format:      document.Format,
// 		Description: document.Description,
// 	}, nil
// }

// func (r *Repository) Update(ctx context.Context, id string, deck domain.Deck) (domain.Deck, error) {
// 	objectID, err := bson.ObjectIDFromHex(id)
// 	if err != nil {
// 		return domain.Deck{}, err
// 	}

// 	document := deckVersionDocument{
// 		Name:        deck.Name,
// 		Format:      deck.Format,
// 		Description: deck.Description,
// 	}

// 	var updatedDocument deckVersionDocument
// 	err = r.collection.FindOneAndUpdate(
// 		ctx,
// 		bson.M{"_id": objectID},
// 		bson.M{"$set": document},
// 		options.FindOneAndUpdate().SetReturnDocument(options.After),
// 	).Decode(&updatedDocument)
// 	if err != nil {
// 		return domain.Deck{}, err
// 	}

// 	return domain.Deck{
// 		ID:          updatedDocument.ID.Hex(),
// 		Name:        updatedDocument.Name,
// 		Format:      updatedDocument.Format,
// 		Description: updatedDocument.Description,
// 	}, nil
// }

// func (r *Repository) Delete(ctx context.Context, id string) error {
// 	objectID, err := bson.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
// 	if err != nil {
// 		return err
// 	}
// 	if result.DeletedCount == 0 {
// 		return mongo.ErrNoDocuments
// 	}

// 	return nil
// }
