package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type StockHistory struct {
	ID        string  `json:"userId"`
	Timestamp uint    `json:"timestamp"`
	Amount    float64 `json:"amount"`
	Price     float64 `json:"price"`
}

type StockCollection struct {
	ID         primitive.ObjectID `json:"_id"`
	StockImage string             `json:"stockImage"`
	Name       string             `json:"name"`
	Sign       string             `json:"sign"`
	Price      float64            `json:"price"`
	History    []StockHistory     `json:"stockHistory"`
}

type TopStock struct {
	Sign  string `json:"sign"`
	Price uint   `json:"price"`
}

type AllStock struct {
	Id    string  `json:"id"`
	Sign  string  `json:"sign"`
	Price float64 `json:"price"`
}

type StockCollectionResponse struct {
	StockImage string         `json:"stockImage"`
	Name       string         `json:"name"`
	Sign       string         `json:"sign"`
	Price      float64        `json:"price"`
	History    []StockHistory `json:"stockHistory"`
}
