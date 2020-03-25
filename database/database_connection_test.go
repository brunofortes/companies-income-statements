package database

import (
	"context"
	"testing"
	"time"
)

type TestEntity struct {
	Name string `bson:"name" json:"name"`
}

func TestDatabaseConnection(t *testing.T) {
	dbConnection := NewDBConnectionWithConfigFile("./config_test.toml")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)

	defer cancel()

	t.Run("must connect to database", func(t *testing.T) {
		_, error := dbConnection.Connect(ctx)

		if error != nil {
			t.Errorf("Error when connecting to database %s", error)
		}
	})

	t.Run("must return a collection", func(t *testing.T) {
		collection, error := dbConnection.GetCollection(ctx, "testEntity")

		if collection == nil {
			t.Errorf("Unable to retrive a collection from database")
		}

		if error != nil {
			t.Errorf("Error when connecting to database %s", error)
		}
	})
}
