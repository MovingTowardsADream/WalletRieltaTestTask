package entity

type Wallet struct {
	ID      string `json:"id"       example:"5b53700ed469fa6a09ea72bb78f36fd9" description:"Уникальный ID кошелька" validate:"required"` //nolint:lll,tagalign // вот так то лучше
	Balance uint   `json:"balance"  example:"100"                              description:"Баланс кошелька"        validate:"required"` //nolint:lll,tagalign // вот так то лучше
}
