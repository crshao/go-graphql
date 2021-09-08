package database

import (
	"context"
	"log"
	"time"

	"github.com/crshao/go-graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	return &DB{
		client: client,
	}
}

func (db *DB) Save(input *model.NewStudent) *model.Student {
	collection := db.client.Database("people").Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, input)

	if err != nil {
		log.Fatal(err)
	}

	return &model.Student{
		ID:   res.InsertedID.(primitive.ObjectID).Hex(),
		Name: input.Name,
		Nim:  input.Nim,
	}
}

func (db *DB) FindByID(ID string) *model.Student {
	ObjectID, err := primitive.ObjectIDFromHex(ID)

	if err != nil {
		log.Fatal(err)
	}

	collection := db.client.Database("people").Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res := collection.FindOne(ctx, bson.M{"_id": ObjectID})

	student := model.Student{}
	res.Decode(&student)
	return &student
}

func (db *DB) All() []*model.Student {
	collection := db.client.Database("people").Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})

	if err != nil {
		log.Fatal(err)
	}

	var students []*model.Student

	for cur.Next(ctx) {
		var student *model.Student
		err := cur.Decode(&student)

		if err != nil {
			log.Fatal(err)
		}
		students = append(students, student)
	}

	return students
}
