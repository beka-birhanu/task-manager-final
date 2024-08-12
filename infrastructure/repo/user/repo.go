/*
Package userrepo provides methods for managing user models in a MongoDB collection.

It supports saving, retrieving by ID or username, and counting users. Errors related to
user operations are handled using custom domain-specific errors.

Dependencies:
- go.mongodb.org/mongo-driver/mongo: MongoDB driver for Go.
- github.com/google/uuid: UUID generation for user IDs.
- github.com/beka-birhanu/domain/errors: Custom domain errors.
- github.com/beka-birhanu/domain/models/user: User model definitions.
*/
package userrepo

import (
	"context"
	"time"

	irepo "github.com/beka-birhanu/task_manager_final/app/common/i_repo"
	errdmn "github.com/beka-birhanu/task_manager_final/domain/errors"
	usermodel "github.com/beka-birhanu/task_manager_final/domain/models/user"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repo handles the persistence of user models.
type Repo struct {
	collection *mongo.Collection
}

// Ensure Repo implements irepo.User.
var _ irepo.User = &Repo{}

// NewRepo creates a new UserRepo with the given MongoDB client, database name, and collection name.
func NewRepo(client *mongo.Client, dbName, collectionName string) *Repo {
	collection := client.Database(dbName).Collection(collectionName)
	return &Repo{
		collection: collection,
	}
}

// Save inserts or updates a user in the repository.
// If the user already exists, it updates the existing record.
// If the user does not exist, it adds a new record.
func (u *Repo) Save(user *usermodel.User) error {
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
			return errdmn.UsernameConflict
		}
		return errdmn.NewUnexpected(err.Error())
	}

	return nil
}

// ById retrieves a user by their ID.
// Returns an error if the user is not found or if an unexpected error occurs.
func (u *Repo) ById(id uuid.UUID) (*usermodel.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	var userBSON usermodel.UserBSON
	if err := u.collection.FindOne(ctx, filter).Decode(&userBSON); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errdmn.UserNotFound
		}
		return nil, errdmn.NewUnexpected(err.Error())
	}
	return usermodel.FromBSON(&userBSON), nil
}

// ByUsername retrieves a user by their username.
// Returns an error if the user is not found or if an unexpected error occurs.
func (u *Repo) ByUsername(username string) (*usermodel.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"username": username}
	var userBSON usermodel.UserBSON
	if err := u.collection.FindOne(ctx, filter).Decode(&userBSON); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errdmn.UserNotFound
		}
		return nil, errdmn.NewUnexpected(err.Error())
	}
	return usermodel.FromBSON(&userBSON), nil
}

// Count returns the total number of users in the repository.
func (u *Repo) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := u.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, errdmn.NewUnexpected(err.Error())
	}
	return count, nil
}
