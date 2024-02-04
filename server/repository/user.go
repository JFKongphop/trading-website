package repository

type UserAccount struct {
	Name         string             `json:"name"`
	ProfileImage string             `json:"profileImage"`
}

type CreateAccount struct {
	Name         string             `json:"name"`
	ProfileImage string             `json:"profileImage"`
	Email        string             `json:"email"`
}

type UserHistory struct {
	Timestamp uint   `json:"timestamp"`
	Name      string `json:"name"`
	Price     uint   `json:"price"`
	Amount    uint   `json:"amount"`
	Status    string `json:"status"`
}


type UserRepository interface {
	Create(CreateAccount) (string, error)
	GetAccount(string) (UserAccount, error)
	GetAllHistory(string) ([]UserHistory, error)
	GetStockHistory(string) ([]UserHistory, error)
	DeleteAccount(string, string) (string, error)
}