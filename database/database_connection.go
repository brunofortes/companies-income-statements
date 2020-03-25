package database

import (
	"context"

	. "github.com/brunofortes/b3-companies-income-statements/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseConnection struct {
	config     Config
	configFile string
	client     *mongo.Client
}

func NewDBConnectionWithConfigFile(configFile string) DatabaseConnection {
	return DatabaseConnection{configFile: configFile}
}

func NewDBConnection() DatabaseConnection {
	return DatabaseConnection{configFile: "./config.toml"}
}

func (d *DatabaseConnection) Connect(ctx context.Context) (*mongo.Client, error) {
	d.config.Read(d.configFile)

	client, error := mongo.NewClient(options.Client().ApplyURI(d.config.Uri))
	if error == nil {
		error = client.Connect(ctx)
	}

	return client, error
}

func (d *DatabaseConnection) Disconnect(ctx context.Context) error {
	return d.client.Disconnect(ctx)
}

func (d *DatabaseConnection) GetCollection(ctx context.Context, collection string) (*mongo.Collection, error) {
	client, error := d.Connect(ctx)
	d.client = client
	return client.Database(d.config.Database).Collection(collection), error
}
