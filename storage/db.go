package db

import "go.mongodb.org/mongo-driver/bson/primitive"

func ToObjectID(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	return oid, err
}
