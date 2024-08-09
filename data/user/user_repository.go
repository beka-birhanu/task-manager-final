package userrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/beka-birhanu/common"
	usermodel "github.com/beka-birhanu/models/user"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserRepo handles the persistence of user models.
type UserRepo struct {
	collection *mongo.Collection
}

// NewUserRepo creates a new UserRepo with the given MongoDB client, database name, and collection name.
func NewUserRepo(client *mongo.Client, dbName, collectionName string) *UserRepo {
	collection := client.Database(dbName).Collection(collectionName)
	return &UserRepo{
		collection: collection,
	}
}

// Save inserts or updates a user in the repository.
// If the user already exists, it updates the existing record.
// If the user does not exist, it adds a new record.
func (u *UserRepo) Save(user *usermodel.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": user.ID()}
	update := bson.M{
		"$set": bson.M{
			"username":     user.Username(),
			"passwordHash": user.PasswordHash(),
			"isAdmin":      user.IsAdmin(),
			"updatedAt":    time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := u.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("conflict: username '%s' already exists", user.Username())
		}
		return err
	}

	return nil
}

// ById retrieves a user by their ID.
func (u *UserRepo) ById(id uuid.UUID) (*usermodel.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	var userBSON usermodel.UserBSON
	if err := u.collection.FindOne(ctx, filter).Decode(&userBSON); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	return usermodel.FromBSON(&userBSON), nil
}

// ByUsername retrieves a user by their username.
func (u *UserRepo) ByUsername(username string) (*usermodel.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"username": username}
	var userBSON usermodel.UserBSON
	if err := u.collection.FindOne(ctx, filter).Decode(&userBSON); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	return usermodel.FromBSON(&userBSON), nil
}
func (u *UserRepo) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return u.collection.CountDocuments(ctx, bson.D{})
}
