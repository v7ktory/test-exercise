package mongodb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/v7ktory/test/internal/config"
	"go.mongodb.org/mongo-driver/bson/mgocompat"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var errNoDBHosts = errors.New("no dbHosts")

type Provider struct {
	DB           *mongo.Database
	QueryTimeout time.Duration
}

func NewMongoDB(ctx context.Context, mongoCfg config.MongoCfg) (*Provider, error) {
	if len(mongoCfg.Hosts) == 0 {
		return nil, errNoDBHosts
	}
	connOpt := options.Client()
	uri := "mongodb://" + strings.Join(mongoCfg.Hosts, ",")
	connOpt.ApplyURI(uri)
	connOpt.SetRegistry(mgocompat.Registry)

	if mongoCfg.Username != "" && mongoCfg.Password != "" {
		connOpt.SetAuth(options.Credential{
			AuthSource: mongoCfg.DB,
			Username:   mongoCfg.Username,
			Password:   mongoCfg.Password,
		})
	}

	queryTimeout := time.Duration(mongoCfg.QueryTimeout)

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, connOpt)
	if err != nil {
		return nil, fmt.Errorf("can't connect to mongodb: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("can't ping mongodb: %w", err)
	}

	return &Provider{
		QueryTimeout: queryTimeout,
		DB:           client.Database(mongoCfg.DB),
	}, nil
}

func (p *Provider) GetCollection(name string) *mongo.Collection {
	return p.DB.Collection(name)
}

func (p *Provider) GetClient() *mongo.Client {
	return p.DB.Client()
}
