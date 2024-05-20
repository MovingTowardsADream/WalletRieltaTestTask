package usecase

import "time"

type Option func(*WalletUseCase)

func Timeout(timeout time.Duration) Option {
	return func(uc *WalletUseCase) {
		uc.timeout = timeout
	}
}

func DefaultBalance(balance uint) Option {
	return func(uc *WalletUseCase) {
		uc.defaultBalance = balance
	}
}
