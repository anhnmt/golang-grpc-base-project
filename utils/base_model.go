package utils

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ IBaseModel = &BaseModel{}

// IBaseModel is the interface that must be implemented by a BaseModel.
type IBaseModel interface {
	CollectionName() string
	PreCreate()
	PreUpdate()
	GetIndexModels() []mongo.IndexModel
}

// BaseModel is a base models struct.
type BaseModel struct {
	Id        string    `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt time.Time `json:"-" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"-" bson:"updated_at,omitempty"`
	DeletedAt time.Time `json:"-" bson:"deleted_at,omitempty"`
}

// CollectionName returns the name of the collection from struct name
func (m *BaseModel) CollectionName() string {
	return CollectionName(m)
}

// PreCreate is a callback that gets called before creating a models.
func (m *BaseModel) PreCreate() {
	m.Id = primitive.NewObjectID().Hex()
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
}

// PreUpdate is a callback that gets called before updating a models.
func (m *BaseModel) PreUpdate() {
	m.UpdatedAt = time.Now()
}

// GetIndexModels returns the index models
func (m *BaseModel) GetIndexModels() []mongo.IndexModel {
	return []mongo.IndexModel{}
}
