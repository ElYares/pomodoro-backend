package repository

import (
	"context"
	"pomodoro-backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoCycleRepository implementa CycleRepository usando MongoDB.
type MongoCycleRepository struct {
	collection *mongo.Collection
}

// NewMongoCycleRepository crea una instancia del repositorio.
func NewMongoCycleRepository(db *mongo.Database) *MongoCycleRepository {
	return &MongoCycleRepository{
		collection: db.Collection("cycles"),
	}
}

// Save guarda un ciclo completado.
func (r *MongoCycleRepository) Save(cycle *domain.PomodoroCycle) error {
	_, err := r.collection.InsertOne(context.Background(), cycle)
	return err
}

// GetByTask obtiene todos los ciclos de una tarea específica.
func (r *MongoCycleRepository) GetByTask(taskID string) ([]*domain.PomodoroCycle, error) {
	filter := bson.M{"task_id": taskID}

	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	var cycles []*domain.PomodoroCycle
	if err := cursor.All(context.Background(), &cycles); err != nil {
		return nil, err
	}

	return cycles, nil
}
