package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type StockHistory struct {
	ID        string              `bson:"userId,omitempty" json:"userId"`
	Timestamp primitive.Timestamp `bson:"timestamp" json:"timestamp"`
	Amount    float64             `bson:"amount" json:"amount"`
	Price     float64             `bson:"price" json:"price"`
}

type StockCollection struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty"`
	StockImage  string              `bson:"stockImage"`
	Name        string              `bson:"name"`
	Sign        string              `bson:"sign"`
	Price       uint                `bson:"price"`
	CreatedDate primitive.Timestamp `bson:"createdDate" json:"createdDate"`
	History     []StockHistory      `bson:"stockHistory"`
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
