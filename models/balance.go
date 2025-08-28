package models

type Balance struct {
	GormModel
	Amount float64 `json:"amount"`
	UserID uint    `json:"user_id" gorm:"column:user_id"`
}

func (Balance) TableName() string {
	return "balances"
}
