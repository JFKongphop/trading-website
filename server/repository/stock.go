package repository

type StockHistory struct {
	ID        string `json:"userId"`
	Timestamp uint   `json:"timestamp"`
	Price     uint   `json:"price"`
}

type StockCollection struct {
	StockImage string         `json:"stockImage"`
	Name       string         `json:"name"`
	Sign       string         `json:"sign"`
	Price      uint           `json:"price"`
	History    []StockHistory `json:"stockHistory"`
}

type TopStock struct {
	Sign  string `json:"sign"`
	Price uint   `json:"price"`
}

type AllStock struct {
	Id    string `json:"id"`
	Sign  string `json:"sign"`
	Price uint   `json:"price"`
}

type StockRepository interface {
	CreateStock(StockCollection) (string, error)
	GetAllStock() ([]AllStock, error)
	GetTopStock() ([]TopStock, error)
	GetStock(string) (StockCollection, error)
	EditStock(string) (string, error)
	DeleteStock(string) (string, error)
}
