package repository

import (
	"context"
	"time"

	"pomodoro-backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTaskRepository struct {
	col *mongo.Collection
}

func NewMongoTaskRepository(db *mongo.Database) *MongoTaskRepository {
	return &MongoTaskRepository{
		col: db.Collection("tasks"),
	}
}

// -----------------------------
// Mongo DTO
// -----------------------------

type mongoTask struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      string             `bson:"user_id"`
	Title       string             `bson:"title"`
	Description string             `bson:"description,omitempty"`
	ProjectID   *string            `bson:"project_id,omitempty"`

	Status      string     `bson:"status"`
	Completed   bool       `bson:"completed"`
	CompletedAt *time.Time `bson:"completed_at,omitempty"`

	PomodorosCompleted int `bson:"pomodoros_completed"`
	TotalFocusMinutes  int `bson:"total_focus_minutes"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

// -----------------------------
// CREATE
// -----------------------------

func (r *MongoTaskRepository) Create(t *domain.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := domainToMongoTask(t)

	res, err := r.col.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		t.ID = oid.Hex()
	}

	return nil
}

// -----------------------------
// UPDATE
// -----------------------------

func (r *MongoTaskRepository) Update(t *domain.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(t.ID)
	if err != nil {
		return err
	}

	doc := domainToMongoTask(t)
	_, err = r.col.ReplaceOne(ctx, bson.M{"_id": oid}, doc)
	return err
}

// -----------------------------
// DELETE
// -----------------------------

func (r *MongoTaskRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

// -----------------------------
// FIND BY ID
// -----------------------------

func (r *MongoTaskRepository) FindByID(id string) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var doc mongoTask
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)
	if err != nil {
		return nil, err
	}

	return mongoToDomainTask(&doc), nil
}

// -----------------------------
// FIND BY USER
// -----------------------------

func (r *MongoTaskRepository) FindByUser(userID string) ([]*domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.col.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*domain.Task
	for cursor.Next(ctx) {
		var doc mongoTask
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		tasks = append(tasks, mongoToDomainTask(&doc))
	}

	return tasks, nil
}

// -----------------------------
// MÉTODOS NUEVOS
// -----------------------------

func (r *MongoTaskRepository) UpdateStatus(id string, status domain.TaskStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.col.UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{
			"status":     string(status),
			"updated_at": time.Now(),
		}},
	)
	return err
}

func (r *MongoTaskRepository) AddRealMinutes(id string, minutes int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.col.UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.M{"$inc": bson.M{
			"total_focus_minutes": minutes,
		}},
	)
	return err
}

func (r *MongoTaskRepository) IncrementPomodoroCount(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.col.UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.M{"$inc": bson.M{
			"pomodoros_completed": 1,
		}},
	)
	return err
}

// -----------------------------
// MAPPERS
// -----------------------------

func domainToMongoTask(t *domain.Task) *mongoTask {
	return &mongoTask{
		UserID:             t.UserID,
		Title:              t.Title,
		Description:        t.Description,
		ProjectID:          t.ProjectID,
		Status:             string(t.Status),
		Completed:          t.Completed,
		CompletedAt:        t.CompletedAt,
		PomodorosCompleted: t.PomodorosCompleted,
		TotalFocusMinutes:  t.TotalFocusMinutes,
		CreatedAt:          t.CreatedAt,
		UpdatedAt:          t.UpdatedAt,
	}
}

func mongoToDomainTask(m *mongoTask) *domain.Task {
	id := ""
	if !m.ID.IsZero() {
		id = m.ID.Hex()
	}

	return &domain.Task{
		ID:                 id,
		UserID:             m.UserID,
		Title:              m.Title,
		Description:        m.Description,
		ProjectID:          m.ProjectID,
		Status:             domain.TaskStatus(m.Status),
		Completed:          m.Completed,
		CompletedAt:        m.CompletedAt,
		PomodorosCompleted: m.PomodorosCompleted,
		TotalFocusMinutes:  m.TotalFocusMinutes,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}
