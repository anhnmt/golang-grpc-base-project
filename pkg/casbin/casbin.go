package casbin

import (
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/xdorro/golang-grpc-base-project/pkg/repo"
)

var _ ICasbin = (*Casbin)(nil)

// ICasbin is the interface that must be implemented by a casbin.
type ICasbin interface {
	Client() *casbin.CachedEnforcer
}

// Casbin is a casbin struct.
type Casbin struct {
	mu          sync.Mutex
	dbName      string
	casbinModel string
	casbinName  string
	client      *casbin.CachedEnforcer

	// options
	repo *repo.Repo
}

// NewCasbin creates a new casbin.
func NewCasbin(repo *repo.Repo) *Casbin {
	c := &Casbin{
		dbName:      viper.GetString("database.name"),
		casbinModel: viper.GetString("casbin.model"),
		casbinName:  viper.GetString("casbin.name"),
		repo:        repo,
	}

	log.Info().Msg("Connecting to Casbin")

	config := &mongodbadapter.AdapterConfig{
		DatabaseName:   c.dbName,
		CollectionName: c.casbinName,
	}
	adapter, err := mongodbadapter.NewAdapterByDB(c.repo.Client(), config)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to create mongodb adapter")
	}

	m, _ := model.NewModelFromString(c.casbinModel)

	enforcer, err := casbin.NewCachedEnforcer(m, adapter)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to create casbin enforcer")
	}

	// Load the policy from DB.
	if err = enforcer.LoadPolicy(); err != nil {
		log.Panic().Err(err).Msg("Failed to load policy")
	}

	// Add enforcer to Casbin.
	c.setClient(enforcer)

	log.Info().Msg("Loaded casbin successfully")

	return c
}

// Client adds a new client to the repository.
func (c *Casbin) setClient(enforcer *casbin.CachedEnforcer) {
	c.mu.Lock()
	c.client = enforcer
	c.mu.Unlock()
}

// Client return cache enforcer
func (c *Casbin) Client() *casbin.CachedEnforcer {
	return c.client
}
