package repository

import (
	"context"
	"time"

	"pomodoro-backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoSessionRepository implementa SessionRepository utilizando MongoDB
// como tecnología de almacenamiento.
type MongoSessionRepository struct {
	col *mongo.Collection
}

// NewMongoSessionRepository construye un nuevo repositorio de sesiones
// acoplado a la colección "sessions" de la base de datos indicada.
func NewMongoSessionRepository(db *mongo.Database) *MongoSessionRepository {
	return &MongoSessionRepository{
		col: db.Collection("sessions"),
	}
}

// mongoSession representa la forma en que se almacenan las sesiones
// físicamente en MongoDB.
type mongoSession struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UserID        string             `bson:"user_id"`
	ProjectID     *string            `bson:"project_id,omitempty"`
	TaskID        *string            `bson:"task_id,omitempty"`
	FocusMinutes  int                `bson:"focus_minutes"`
	BreakMinutes  int                `bson:"break_minutes"`
	State         string             `bson:"state"`
	StartedAt     time.Time          `bson:"started_at"`
	PausedAt      *time.Time         `bson:"paused_at,omitempty"`
	FinishedAt    *time.Time         `bson:"finished_at,omitempty"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
	Interruptions int                `bson:"interruptions"`
}

// CreateSession inserta una nueva sesión en la colección de MongoDB.
func (r *MongoSessionRepository) CreateSession(s *domain.Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := domainToMongo(s)

	res, err := r.col.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		s.ID = oid.Hex()
	}

	return nil
}

// UpdateSession reemplaza el documento asociado a la sesión por su nueva versión.
func (r *MongoSessionRepository) UpdateSession(s *domain.Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(s.ID)
	if err != nil {
		return err
	}

	doc := domainToMongo(s)

	_, err = r.col.ReplaceOne(ctx, bson.M{"_id": oid}, doc)
	return err
}

// FindByID recupera una sesión con base en su identificador único.
func (r *MongoSessionRepository) FindByID(id string) (*domain.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var doc mongoSession
	if err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc); err != nil {
		return nil, err
	}

	return mongoToDomain(&doc), nil
}

// domainToMongo proyecta una entidad de dominio hacia su representación
// específica para MongoDB.
func domainToMongo(s *domain.Session) *mongoSession {
	return &mongoSession{
		UserID:        s.UserID,
		ProjectID:     s.ProjectID,
		TaskID:        s.TaskID,
		FocusMinutes:  s.FocusMinutes,
		BreakMinutes:  s.BreakMinutes,
		State:         string(s.State),
		StartedAt:     s.StartedAt,
		PausedAt:      s.PausedAt,
		FinishedAt:    s.FinishedAt,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
		Interruptions: s.Interruptions,
	}
}

// mongoToDomain proyecta un documento de MongoDB hacia la entidad de dominio.
func mongoToDomain(m *mongoSession) *domain.Session {
	id := ""
	if !m.ID.IsZero() {
		id = m.ID.Hex()
	}

	return &domain.Session{
		ID:            id,
		UserID:        m.UserID,
		ProjectID:     m.ProjectID,
		TaskID:        m.TaskID,
		FocusMinutes:  m.FocusMinutes,
		BreakMinutes:  m.BreakMinutes,
		State:         domain.SessionState(m.State),
		StartedAt:     m.StartedAt,
		PausedAt:      m.PausedAt,
		FinishedAt:    m.FinishedAt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		Interruptions: m.Interruptions,
	}
}
