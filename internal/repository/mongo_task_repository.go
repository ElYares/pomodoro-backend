package repository

import (
	"context"
	"time"

	"pomodoro-backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoTaskRepository
//
// Implementación concreta del repositorio de tareas utilizando MongoDB como
// mecanismo de persistencia. Esta capa traduce estructuras del dominio a
// documentos BSON y viceversa.
//
// Se mantiene completamente aislada del dominio mediante DTOs internos
// (mongoTask), de forma que cambios en la base de datos no afecten al dominio.
type MongoTaskRepository struct {
	col *mongo.Collection
}

// NewMongoTaskRepository inicializa un repositorio sobre la colección "tasks".
func NewMongoTaskRepository(db *mongo.Database) *MongoTaskRepository {
	return &MongoTaskRepository{
		col: db.Collection("tasks"),
	}
}

// mongoTask
//
// Estructura interna utilizada para serializar/deserializar documentos BSON
// desde y hacia MongoDB. Separamos explícitamente este esquema del dominio
// para evitar acoplamiento directo.
type mongoTask struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      string             `bson:"user_id"`
	Title       string             `bson:"title"`
	Description string             `bson:"description,omitempty"`
	ProjectID   *string            `bson:"project_id,omitempty"`
	Completed   bool               `bson:"completed"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
	CompletedAt *time.Time         `bson:"completed_at,omitempty"`
}

// Create inserta una nueva tarea en la colección de MongoDB.
func (r *MongoTaskRepository) Create(t *domain.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := domainToMongoTask(t)

	res, err := r.col.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	// El ID generado por Mongo se asigna nuevamente al modelo de dominio.
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		t.ID = oid.Hex()
	}

	return nil
}

// Update reemplaza el documento asociado utilizando un ReplaceOne completo.
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

// Delete elimina una tarea directamente por su identificador.
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

// FindByID devuelve una tarea específica si existe en la colección.
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

// FindByUser obtiene todas las tareas asociadas a un usuario dado.
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

// Mapping helpers
// Convierten entre el modelo de dominio y el modelo de persistencia.
func domainToMongoTask(t *domain.Task) *mongoTask {
	return &mongoTask{
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		ProjectID:   t.ProjectID,
		Completed:   t.Completed,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		CompletedAt: t.CompletedAt,
	}
}

func mongoToDomainTask(m *mongoTask) *domain.Task {
	id := ""
	if !m.ID.IsZero() {
		id = m.ID.Hex()
	}

	return &domain.Task{
		ID:          id,
		UserID:      m.UserID,
		Title:       m.Title,
		Description: m.Description,
		ProjectID:   m.ProjectID,
		Completed:   m.Completed,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		CompletedAt: m.CompletedAt,
	}
}
