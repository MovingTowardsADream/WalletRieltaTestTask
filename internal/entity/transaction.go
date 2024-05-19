package entity

import "time"

type Transaction struct {
	Time   time.Time `json:"time"   example:"2024-02-04T17:25:35.448Z"         description:"Дата и время перевода"  validate:"required" format:"date-time"`  //nolint:lll,tagalign // вот так то лучше
	From   string    `json:"from"   example:"5b53700ed469fa6a09ea72bb78f36fd9" description:"ID исходящего кошелька" validate:"required" pg:"from_wallet_id"` //nolint:lll,tagalign // вот так то лучше
	To     string    `json:"to"     example:"eb376add88bf8e70f80787266a0801d5" description:"ID входящего кошелька"  validate:"required" pg:"to_wallet_id"`   //nolint:lll,tagalign // вот так то лучше
	Amount uint      `json:"amount" example:"30"                               description:"Сумма перевода"         validate:"required"`                     //nolint:lll,tagalign // вот так то лучше
}
