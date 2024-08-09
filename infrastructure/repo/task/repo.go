/*
Package taskrepo provides methods for managing tasks in a MongoDB collection.

It supports adding, updating, deleting, and retrieving tasks. Errors related to task
operations are handled using custom domain-specific errors.

Dependencies:
- go.mongodb.org/mongo-driver/mongo: MongoDB driver for Go.
- github.com/google/uuid: UUID generation for task IDs.
- github.com/beka-birhanu/domain/errors: Custom domain errors.
- github.com/beka-birhanu/domain/models/task: Task model definitions.
*/
package taskrepo

import (
	"context"
	"time"

	irepo "github.com/beka-birhanu/app/common/i_repo"
	errdmn "github.com/beka-birhanu/domain/errors"
	taskmodel "github.com/beka-birhanu/domain/models/task"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repo represents a repository for managing tasks.
type Repo struct {
	collection *mongo.Collection
}

// Ensure Repo implements irepo.Task
var _ irepo.Task = &Repo{}

// New creates a new Repo for managing tasks with the given MongoDB client, database name, and collection name.
func New(client *mongo.Client, dbName, collectionName string) *Repo {
	collection := client.Database(dbName).Collection(collectionName)
	return &Repo{
		collection: collection,
	}
}

// createScopedContext creates a new context with a timeout for scoped operations.
func createScopedContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

// Add adds a new task to the collection. Returns an error if there is an ID conflict.
func (r *Repo) Add(task *taskmodel.Task) error {
	ctx, cancel := createScopedContext()
	defer cancel()

	taskBSON := task.ToBSON()

	// Retry task creation if there's an ID conflict
	for {
		_, err := r.collection.InsertOne(ctx, taskBSON)
		if mongo.IsDuplicateKeyError(err) {
			continue
		} else if err != nil {
			return errdmn.NewUnexpected(err.Error())
		}
		return nil
	}
}

// Update updates an existing task. Returns an error if the task is not found.
func (r *Repo) Update(task *taskmodel.Task) error {
	ctx, cancel := createScopedContext()
	defer cancel()

	filter := bson.M{"_id": task.ID}
	update := bson.M{
		"$set": bson.M{
			"title":       task.Title,
			"description": task.Description,
			"dueDate":     task.DueDate,
			"status":      task.Status,
			"updatedAt":   time.Now(),
		},
	}

	result := r.collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return errdmn.TaskNotFound
		}
		return errdmn.NewUnexpected(result.Err().Error())
	}
	return nil
}

// Delete removes a task by ID. Returns an error if the task is not found.
func (r *Repo) Delete(id uuid.UUID) error {
	ctx, cancel := createScopedContext()
	defer cancel()

	filter := bson.M{"_id": id}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return errdmn.NewUnexpected(err.Error())
	}
	if result.DeletedCount == 0 {
		return errdmn.TaskNotFound
	}
	return nil
}

// GetAll returns a list of all tasks.
func (r *Repo) GetAll() ([]*taskmodel.Task, error) {
	ctx, cancel := createScopedContext()
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, errdmn.NewUnexpected(err.Error())
	}
	defer cursor.Close(ctx)

	var tasks []*taskmodel.Task
	for cursor.Next(ctx) {
		var taskBSON taskmodel.TaskBSON
		if err := cursor.Decode(&taskBSON); err != nil {
			return nil, errdmn.NewUnexpected(err.Error())
		}
		tasks = append(tasks, taskmodel.FromBSON(&taskBSON))
	}
	if err := cursor.Err(); err != nil {
		return nil, errdmn.NewUnexpected(err.Error())
	}
	return tasks, nil
}

// GetSingle returns a task by ID. Returns an error if the task is not found.
func (r *Repo) GetSingle(id uuid.UUID) (*taskmodel.Task, error) {
	ctx, cancel := createScopedContext()
	defer cancel()

	filter := bson.M{"_id": id}
	var taskBSON taskmodel.TaskBSON
	if err := r.collection.FindOne(ctx, filter).Decode(&taskBSON); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errdmn.TaskNotFound
		}
		return nil, errdmn.NewUnexpected(err.Error())
	}
	task := taskmodel.FromBSON(&taskBSON)
	return task, nil
}
