package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserStock struct {
	Name   string `bson:"name" json:"name"`
	Amount uint   `bson:"amount" json:"amount"`
}

type UserAccount struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name         string             `bson:"name" json:"name"`
	ProfileImage string             `bson:"profileImage" json:"profileImage"`
	Email        string             `bson:"email" json:"email"`
	RegisterDate uint               `bson:"registerDate" json:"registerDate"`
	History      []UserHistory      `bson:"userHistory" json:"userHistory"`
	Stock        []UserStock        `bson:"userStock" json:"userStock"`
}

type CreateAccount struct {
	Name         string             `json:"name"`
	ProfileImage string             `json:"profileImage"`
	Email        string             `json:"email"`
}

type UserHistory struct {
	Timestamp uint   `bson:"timestamp" json:"timestamp"`
	Name      string `bson:"name" json:"name"`
	Price     uint   `bson:"price" json:"price"`
	Amount    uint   `bson:"amount" json:"amount"`
	Status    string `bson:"status" json:"status"`
}

type UserResponse struct {
	Name         string             `json:"name"`
	ProfileImage string             `json:"profileImage"`
}