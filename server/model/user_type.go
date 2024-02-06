package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStock struct {
	StockId string  `bson:"stockId" json:"stockId"`
	Amount  float64 `bson:"amount" json:"amount"`
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
	Name         string `json:"name"`
	ProfileImage string `json:"profileImage"`
	Email        string `json:"email"`
}

type OrderRequest struct {
	StockId     string  `json:"stockId"`
	UserId      string  `json:"userId"`
	Price       float64 `json:"price"`
	Amount      float64 `json:"amount"`
	OrderType   string  `json:"orderType"`   // auto, order
	OrderMethod string  `json:"orderMethod"` // buy, sale
}

type UserHistory struct {
	Timestamp   uint    `bson:"timestamp" json:"timestamp"`
	StockId     string  `bson:"stockId" json:"stockId"`
	Price       float64 `bson:"price" json:"price"`
	Amount      float64 `bson:"amount" json:"amount"`
	Status      string  `bson:"status" json:"status"`           // pending, success, cancle
	OrderType   string  `bson:"orderType" json:"orderType"`     // auto, order
	OrderMethod string  `bson:"orderMethod" json:"orderMethod"` // buy, sale
}

type UserResponse struct {
	Name         string `json:"name"`
	ProfileImage string `json:"profileImage"`
}

type OrderOperation struct {
	DB          *mongo.Collection  `json:"db"`
	UserId      primitive.ObjectID `json:"userId"`
	StockId     string             `json:"stockId"`
	History     UserHistory        `json:"userHistory"`
	Amount      float64            `json:"amount"`
	OrderMethod string             `json:"orderMethod"`
}
