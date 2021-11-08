package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var ctx context.Context
var err error
var sb []byte

type Recipe struct {
	//swagger:ignore
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Tags         []string           `json:"tags" bson:"tags"`
	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
	Instructions []string           `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt" bson:"publishedAt"`
}

var listOfRecipes []Recipe

var recipes []struct {
	ID           string    `json:"id" bson:"_id"`
	Name         string    `json:"name" bson:"name"`
	Tags         []string  `json:"tags" bson:"tags"`
	Ingredients  []string  `json:"ingredients" bson:"ingredients"`
	Instructions []string  `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time `json:"publishedAt" bson:"publishedAt"`
}

func init() {
	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		fmt.Println("Error connecting to mongodb server:", err)
		return
	}

	sb, err = ioutil.ReadFile("recipes.json")
	if err != nil {
		fmt.Println("Error reading recipes.json file:", err)
		return
	}

	if err := json.Unmarshal(sb, &recipes); err != nil {
		fmt.Println("Error unmarshalling from json to recipes struct:", err)
		return
	}
}

func main() {
	collection := client.Database("nginxRecipes").Collection("recipes")

	for _, v := range recipes {
		var recipe Recipe
		recipe.ID = primitive.NewObjectID()
		recipe.Name = v.Name
		recipe.Tags = v.Tags
		recipe.Ingredients = v.Ingredients
		recipe.Instructions = v.Instructions
		recipe.PublishedAt = v.PublishedAt
		collection.InsertOne(ctx, &recipe)
	}
}
