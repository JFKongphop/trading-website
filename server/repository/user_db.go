package repository

import "go.mongodb.org/mongo-driver/mongo"

type userRepositoryDB struct {
	db *mongo.Collection
}

func NewUserRepositoryDB(db *mongo.Collection) UserRepository {
	return userRepositoryDB{db}
}

func (r userRepositoryDB) Create(account CreateAccount) (string, error) {
	return "", nil
}

func (r userRepositoryDB) GetAccount(string) (UserAccount, error) {
	return UserAccount{}, nil
}

func (r userRepositoryDB) GetAllHistory(string) ([]UserHistory, error) {
	return []UserHistory{}, nil
}

func (r userRepositoryDB) GetStockHistory(string) ([]UserHistory, error) {
	return []UserHistory{}, nil
}

func (r userRepositoryDB) DeleteAccount(string, string) (string, error) {
	return "", nil
}
