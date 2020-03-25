package company

import (
	"context"

	. "github.com/brunofortes/b3-companies-income-statements/database"
	. "github.com/brunofortes/b3-companies-income-statements/financial"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	COLLECTION = "companies"
)

type Company struct {
	ID         primitive.ObjectID `json:"ID" bson:"_id,omitempty"`
	Name       string             `bson:"name" json:"name"`
	Labels     []string           `bson:"tickets" json:"tickets"`
	Sectors    []string           `bson:"sectors" json:"sectors"`
	Industries []string           `bson:"industries" json:"industries"`
	Financials []Financial        `bson:"financials" json:"financials"`
}

type CompanyConnection struct {
	databaseConnection  DatabaseConnection
	companiesCollection *mongo.Collection
	ctx                 context.Context
}

func NewCompanyConnection() CompanyConnection {
	ctx := context.Background()
	databaseConnection := NewDBConnection()
	companiesCollection, _ := databaseConnection.GetCollection(ctx, COLLECTION)
	return CompanyConnection{
		databaseConnection:  databaseConnection,
		companiesCollection: companiesCollection,
		ctx:                 ctx,
	}
}

func (c *CompanyConnection) Save(company Company) error {

	if company.ID.IsZero() {
		_, err := c.companiesCollection.InsertOne(c.ctx, &company)
		return err
	} else {
		_, err := c.companiesCollection.ReplaceOne(c.ctx, bson.M{"_id": &company.ID}, &company)
		return err
	}

}

func (c *CompanyConnection) ListAll() ([]Company, error) {
	result := []Company{}

	cur, error := c.companiesCollection.Find(c.ctx, bson.M{})
	if error == nil {
		error = cur.All(c.ctx, &result)
	}

	return result, error
}

func (c *CompanyConnection) FindByName(name string) ([]Company, error) {
	result := []Company{}

	cur, error := c.companiesCollection.Find(c.ctx, bson.M{"name": name})
	if error == nil {
		error = cur.All(c.ctx, &result)
	}

	return result, error
}

func (c *CompanyConnection) Disconnect() {
	c.databaseConnection.Disconnect(c.ctx)
	c.ctx.Done()
}
