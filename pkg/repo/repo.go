package repo

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/xdorro/golang-grpc-base-project/utils"
)

var _ IRepo = (*Repo)(nil)

// IRepo is the interface that must be implemented by a repository.
type IRepo interface {
	Client() *mongo.Client
	Close() error
	Database() *mongo.Database
	Collection(name string) *mongo.Collection
	CollectionModel(model utils.IBaseModel) *mongo.Collection
}

// Repo is a repository struct.
type Repo struct {
	mu     sync.Mutex
	dbURL  string
	dbName string

	client *mongo.Client
}

// NewRepo creates a new repository.
func NewRepo() *Repo {
	r := &Repo{
		dbURL:  viper.GetString("database.url"),
		dbName: viper.GetString("database.name"),
	}

	log.Info().
		Str("db_url", r.dbURL).
		Str("db_name", r.dbName).
		Msg("Connecting to MongoDB")

	// Set connect timeout to 15 seconds
	ctxConn, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().
		ApplyURI(r.dbURL)

	// Create new client and connect to MongoDB
	client, err := mongo.Connect(ctxConn, clientOpts)
	if err != nil {
		log.Panic().Err(err).Msg("Connecting to MongoDB failed")
	}

	// Ping the primary
	if err = client.Ping(ctxConn, nil); err != nil {
		log.Panic().Err(err).Msg("Ping to MongoDB failed")
	}

	// Add client to repository
	r.setClient(client)

	log.Info().Msg("Connecting to MongoDB successfully.")

	return r
}

// Close closes the repository.
func (r *Repo) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.client.Disconnect(ctx); err != nil {
		log.Err(err).Msg("Failed to disconnect from MongoDB")
		return err
	}

	return nil
}

// Database returns the mongo database by Name.
func (r *Repo) Database() *mongo.Database {
	return r.client.Database(r.dbName)
}

// Collection returns the mongo collection by Name.
func (r *Repo) Collection(name string) *mongo.Collection {
	return r.Database().Collection(name)
}

// CollectionModel returns the mongo collection models by Name.
func (r *Repo) CollectionModel(model utils.IBaseModel) *mongo.Collection {
	return r.Collection(model.CollectionName())
}

// Client returns the mongo client
func (r *Repo) Client() *mongo.Client {
	return r.client
}

// setClient adds a new client to the repository.
func (r *Repo) setClient(client *mongo.Client) {
	r.mu.Lock()
	r.client = client
	r.mu.Unlock()
}
