package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Product represents a cloth item in the store
type Product struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Name        string             `bson:"name" json:"name"`
    Description string             `bson:"description" json:"description"`
    Price       float64            `bson:"price" json:"price"`
    Category    string             `bson:"category" json:"category"`
}

