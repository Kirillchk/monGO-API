package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const conectionString = "mongodb://localhost:27017"
const DBName = "MonGO"

var globclient *mongo.Client
var database *mongo.Database

type User struct {
	Login    string
	Password string
	JWT      string
}

func processError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func InitDB() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(conectionString))
	processError(err)
	globclient = client
	database = globclient.Database(DBName)
}

// User interactions
func FindUser(username string) (string, error) {
	collection := database.Collection("users")
	filter := bson.M{"login": username}
	var userInfo User
	res := collection.FindOne(context.TODO(), filter).Decode(&userInfo)
	if res == mongo.ErrNoDocuments {
		return "", mongo.ErrNoDocuments
	} else {
		return userInfo.JWT, nil
	}
}
func AddUser(username string, password string) (string, error) {
	collection := database.Collection("users")

	filter := bson.M{"login": username}
	res := collection.FindOne(context.TODO(), filter).Err()
	if res == mongo.ErrNoDocuments {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": username,                              // subject (user id)
			"iat": time.Now().Unix(),                     // issued at
			"exp": time.Now().Add(time.Hour * 24).Unix(), // expiration time
		})
		signed, _ := token.SignedString([]byte("7f3a1e8d2b6c5f9a0d4e7b2c8f1a6d3e"))
		doc := map[string]interface{}{
			"Login":    username,
			"Password": password,
			"JWT":      signed,
		}
		collection.InsertOne(context.TODO(), doc)
		return signed, nil
	} else {
		return "", errors.New("already exist")
	}
}

// Collections interactions
func ListCollections() []string {
	collection, err := database.ListCollectionNames(context.TODO(), bson.D{})
	processError(err)
	return collection
}
func AddColletion(collectionName string) {
	database.CreateCollection(context.TODO(), collectionName)
}
func FindCollection(collectionName string) ([]bson.M, error) {
	var collectionRaw []bson.M
	collection := database.Collection(collectionName)
	cursor, err := collection.Find(context.TODO(), bson.M{})
	processError(err)
	err = cursor.All(context.TODO(), &collectionRaw)
	processError(err)
	return collectionRaw, err
}
func DeleteCollection(collectionName string) {
	database.Collection(collectionName).Drop(context.TODO())
}

// Document interactions
type Document struct {
	collection string
	document   bson.M
}

func (d *Document) CollRef() *mongo.Collection {
	return database.Collection(d.collection)
}
func (d *Document) Add() {
	d.CollRef().InsertOne(context.TODO(), d.document)
}
func (d *Document) Update(newDocument bson.M) {
	_, err := d.CollRef().UpdateOne(context.TODO(), d.document, newDocument)
	fmt.Println(err)
}
func (d *Document) Delete() {
	d.CollRef().DeleteOne(context.TODO(), d.document)
}
