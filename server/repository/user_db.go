package repository

import (
	"context"
	"server/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ctx = context.Background()

type userRepositoryDB struct {
	db *mongo.Collection
}
type UserStock = model.UserStock

func NewUserRepositoryDB(db *mongo.Collection) UserRepository {
	return userRepositoryDB{db}
}

func (r userRepositoryDB) Create(data CreateAccount) (string, error) {
	user := UserAccount{
		Name: data.Name,
		ProfileImage: data.ProfileImage,
		Email: data.Email,
		RegisterDate: uint(time.Now().Unix()),
		History:      []UserHistory{},
		Stock:        []UserStock{},
	}

	_, err := r.db.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	return "Successfully created account", nil
}

func (r userRepositoryDB) GetAccount(id string) (userAccount UserAccount, err error) {
	filter := bson.M{ "_id": id }
	err = r.db.FindOne(ctx, filter).Decode(&userAccount)
	if err != nil {
		return UserAccount{}, err
	}

	return userAccount, nil
}

func (r userRepositoryDB) GetAllHistories(string) ([]UserHistory, error) {
	return []UserHistory{}, nil
}

func (r userRepositoryDB) GetStockHistory(string) ([]UserHistory, error) {
	return []UserHistory{}, nil
}

func (r userRepositoryDB) DeleteAccount(string, string) (string, error) {
	return "", nil
}
