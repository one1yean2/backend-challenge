package repositories

import (
	"context"
	"one1-be-chal/internal/core/domain"
	"one1-be-chal/internal/core/ports"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database, collectionName string) ports.UserRepository {
	return &MongoUserRepository{
		collection: db.Collection(collectionName),
	}
}

func (u *MongoUserRepository) Save(ctx context.Context, user domain.User) error {
	if _, err := u.collection.InsertOne(ctx, user); err != nil {
		return err
	}
	return nil
}

func (u *MongoUserRepository) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	var user domain.User
	if err := u.collection.FindOne(ctx, bson.M{"id": id}).Decode(&user); err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *MongoUserRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	cursor, err := u.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (u *MongoUserRepository) UpdateUser(ctx context.Context, uid string, user bson.M) error {
	_, err := u.collection.UpdateOne(
		ctx,
		bson.M{"id": uid},
		bson.M{"$set": user},
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *MongoUserRepository) DeleteUser(ctx context.Context, id string) error {
	if _, err := u.collection.DeleteOne(ctx, bson.M{"id": id}); err != nil {
		return err
	}
	return nil
}

func (u *MongoUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user *domain.User
	if err := u.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *MongoUserRepository) GetUserCount(ctx context.Context) (int64, error) {
	count, err := u.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return count, nil
}
