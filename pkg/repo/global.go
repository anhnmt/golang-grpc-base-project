package repo

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Find return an array of objects
func Find[T any](collection *mongo.Collection, filter any, opt ...*options.FindOptions) ([]*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, filter, opt...)
	if err != nil {
		log.Err(err).Msg("Error find all users")
		return nil, err
	}

	defer func() {
		_ = cur.Close(ctx)
	}()

	var data []*T
	for cur.Next(ctx) {
		obj := new(T)
		if err = cur.Decode(obj); err != nil {
			log.Err(err).Msg("Error find all")
			return nil, err
		}

		data = append(data, obj)
	}

	return data, nil
}

// CountDocuments returns the number of documents
func CountDocuments(collection *mongo.Collection, filter any) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Err(err).Msg("Error count all users")
		return 0, err
	}

	return total, nil
}

// FindOne return an object
func FindOne[T any](collection *mongo.Collection, filter any, opt ...*options.FindOneOptions) (
	*T, error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result := new(T)
	err := collection.FindOne(ctx, filter, opt...).Decode(result)
	if err != nil {
		log.Err(err).Msg("Error find user")
		return nil, err
	}

	return result, nil
}

// InsertOne inserts one
func InsertOne(collection *mongo.Collection, data any, opt ...*options.InsertOneOptions) (
	*mongo.InsertOneResult, error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, data, opt...)
	if err != nil {
		log.Err(err).Msg("Error creating user")
		return nil, err
	}

	return res, nil
}

// InsertMany insert many
func InsertMany(collection *mongo.Collection, data []any, opt ...*options.InsertManyOptions) (
	*mongo.InsertManyResult, error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := collection.InsertMany(ctx, data, opt...)
	if err != nil {
		log.Err(err).Msg("Error insert many")
		return nil, err
	}

	return res, nil
}

// UpdateOne updates one
func UpdateOne(collection *mongo.Collection, filter, data any, opt ...*options.UpdateOptions) (
	*mongo.UpdateResult, error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := collection.UpdateOne(ctx, filter, data, opt...)
	if err != nil {
		log.Err(err).Msg("Error updating user")
		return nil, err
	}

	return res, nil
}

// UpdateMany updates many
func UpdateMany(collection *mongo.Collection, filter, data any, opt ...*options.UpdateOptions) (
	*mongo.UpdateResult, error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := collection.UpdateMany(ctx, filter, data, opt...)
	if err != nil {
		log.Err(err).Msg("Error updating user")
		return nil, err
	}

	return res, nil
}

// DeleteOne deletes one
func DeleteOne(collection *mongo.Collection, filter any, opt ...*options.DeleteOptions) (
	*mongo.DeleteResult, error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := collection.DeleteOne(ctx, filter, opt...)
	if err != nil {
		log.Err(err).Msg("Error deleting user")
		return nil, err
	}

	return res, nil
}

// DeleteMany deletes many
func DeleteMany(collection *mongo.Collection, filter any, opt ...*options.DeleteOptions) (
	*mongo.DeleteResult, error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := collection.DeleteMany(ctx, filter, opt...)
	if err != nil {
		log.Err(err).Msg("Error deleting user")
		return nil, err
	}

	return res, nil
}

// SoftDeleteOne soft deletes one
func SoftDeleteOne(collection *mongo.Collection, filter any, opt ...*options.UpdateOptions) (
	*mongo.UpdateResult, error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	data := bson.M{
		"$set": bson.D{
			{
				"deleted_at",
				time.Now(),
			},
		},
	}
	res, err := collection.UpdateOne(ctx, filter, data, opt...)
	if err != nil {
		log.Err(err).Msg("Error soft deleting user")
		return nil, err
	}

	return res, nil
}
