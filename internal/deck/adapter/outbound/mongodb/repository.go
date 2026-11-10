package mongodb

import (
	"context"

	"github.com/ropehapi/mtg-decklogger/internal/deck/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository struct {
	collection *mongo.Collection
}

type deckDocument struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Format      string        `bson:"format"`
	Description string        `bson:"description"`
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

func (repository *Repository) ExistsByID(ctx context.Context, id string) (bool, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	err = repository.collection.FindOne(ctx, bson.M{"_id": objectID}).Err()
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
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
			ID:          document.ID.Hex(),
			Name:        document.Name,
			Format:      document.Format,
			Description: document.Description,
		})
	}

	return decks, nil
}

func (r *Repository) FindById(ctx context.Context, id string) (domain.Deck, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.Deck{}, err
	}

	var document deckDocument
	if err := r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&document); err != nil {
		return domain.Deck{}, err
	}

	return domain.Deck{
		Name:        document.Name,
		Format:      document.Format,
		Description: document.Description,
	}, nil
}

func (r *Repository) Update(ctx context.Context, id string, deck domain.Deck) (domain.Deck, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.Deck{}, err
	}

	document := deckDocument{
		Name:        deck.Name,
		Format:      deck.Format,
		Description: deck.Description,
	}

	var updatedDocument deckDocument
	err = r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": document},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedDocument)
	if err != nil {
		return domain.Deck{}, err
	}

	return domain.Deck{
		ID:          updatedDocument.ID.Hex(),
		Name:        updatedDocument.Name,
		Format:      updatedDocument.Format,
		Description: updatedDocument.Description,
	}, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
