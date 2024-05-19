package entity

type CreateNewWalletWithBalanceRequest struct {
	Balance uint `json:"balance"`
}

type SendFundsRequest struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount uint   `json:"amount"`
}

type GetWalletHistoryByIDRequest struct {
	WalletID string `json:"walletId"`
}

type GetWalletByIDRequest struct {
	WalletID string `json:"walletId"`
}
