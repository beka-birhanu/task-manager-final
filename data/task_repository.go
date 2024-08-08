package data

import (
	"context"
	"time"

	"github.com/beka-birhanu/common"
	"github.com/beka-birhanu/models/taskmodel"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TaskService represents a service for managing tasks.
type TaskService struct {
	collection *mongo.Collection
}

// NewTaskService creates a new TaskService.
func NewTaskService(client *mongo.Client, dbName, collectionName string) *TaskService {
	collection := client.Database(dbName).Collection(collectionName)
	return &TaskService{
		collection: collection,
	}
}

// createScopedContext creates a new context with a timeout for scoped operations.
func createScopedContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

// Add adds a new task to the collection. Returns an error if there is an ID conflict.
func (s *TaskService) Add(title, description, status string, dueDate time.Time) (*taskmodel.Task, error) {
	ctx, cancel := createScopedContext()
	defer cancel()

	taskConfig := taskmodel.TaskConfig{
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Status:      status,
	}

	// Recreate task until the ID conflict is resolved
	for {
		task, err := taskmodel.New(taskConfig)
		if err != nil {
			return nil, err
		}

		taskBSON := task.ToBSON()
		_, err = s.collection.InsertOne(ctx, taskBSON)
		if mongo.IsDuplicateKeyError(err) {
			// If a duplicate key error occurs, generate a new ID and try again
			continue
		} else if err != nil {
			return nil, err
		}
		return task, nil
	}
}

// Update updates an existing task. Returns an error if the task is not found.
func (s *TaskService) Update(id uuid.UUID, title, description, status string, dueDate time.Time) (*taskmodel.Task, error) {
	ctx, cancel := createScopedContext()
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"title":       title,
			"description": description,
			"dueDate":     dueDate,
			"status":      status,
			"updatedAt":   time.Now(),
		},
	}

	result := s.collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, common.ErrNotFound
		}
		return nil, result.Err()
	}

	var taskBSON taskmodel.TaskBSON
	if err := result.Decode(&taskBSON); err != nil {
		return nil, err
	}
	task := taskmodel.FromBSON(&taskBSON)
	return task, nil
}

// Delete removes a task by ID. Returns an error if the task is not found.
func (s *TaskService) Delete(id uuid.UUID) error {
	ctx, cancel := createScopedContext()
	defer cancel()

	filter := bson.M{"_id": id}
	result, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return common.ErrNotFound
	}
	return nil
}

// GetAll returns a list of pointers to all tasks.
func (s *TaskService) GetAll() ([]*taskmodel.Task, error) {
	ctx, cancel := createScopedContext()
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*taskmodel.Task
	for cursor.Next(ctx) {
		var taskBSON taskmodel.TaskBSON
		if err := cursor.Decode(&taskBSON); err != nil {
			return nil, err
		}
		tasks = append(tasks, taskmodel.FromBSON(&taskBSON))
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetSingle returns a task by ID. Returns an error if the task is not found.
func (s *TaskService) GetSingle(id uuid.UUID) (*taskmodel.Task, error) {
	ctx, cancel := createScopedContext()
	defer cancel()

	filter := bson.M{"_id": id}
	var taskBSON taskmodel.TaskBSON
	if err := s.collection.FindOne(ctx, filter).Decode(&taskBSON); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	task := taskmodel.FromBSON(&taskBSON)
	return task, nil
}
