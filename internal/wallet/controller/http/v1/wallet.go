package v1

import (
	"WalletRieltaTestTask/internal/entity"
	"WalletRieltaTestTask/internal/wallet/usecase"
	"WalletRieltaTestTask/pkg/logger"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type walletRoutes struct {
	w usecase.Wallet
	l *slog.Logger
}

func newWalletRoutes(handler *gin.RouterGroup, w usecase.Wallet, l *slog.Logger) {
	r := &walletRoutes{w, l}

	h := handler.Group("/wallet")
	{
		h.POST("", r.createNewWallet)
		h.POST("/:walletId/send", r.sendFunds)
		h.GET("/:walletId/history", r.GetWalletHistoryByID)
		h.GET("/:walletId", r.GetWalletByID)
	}
}

// @Summary     Создание кошелька
// @Description Создает новый кошелек с уникальным ID. Идентификатор генерируется сервером.
// @Description
// @Description Созданный кошелек должен иметь сумму 100.0 у.е. на балансе
// @Tags  	    Wallet
// @Success     200 {object} entity.Wallet "Кошелек создан"
// @Failure     500 "Не удалось создать кошелек"
// @Failure     504 "Время ожидания вышло"
// @Router      /wallet [post].
func (r *walletRoutes) createNewWallet(c *gin.Context) {
	wallet, err := r.w.CreateNewWalletWithDefaultBalance(c.Request.Context())
	if err != nil {
		if errors.Is(err, entity.ErrTimeout) {
			c.AbortWithStatus(http.StatusGatewayTimeout)
			return
		}

		r.l.Error("http - v1 - createNewWallet", logger.Err(err))
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	c.JSON(http.StatusOK, wallet)
}

// @Description Запрос перевода средств.
type transactionRequest struct {
	To     string `json:"to"     example:"eb376add88bf8e70f80787266a0801d5" description:"ID кошелька, куда нужно перевести деньги" validate:"required"` //nolint:lll,tagalign // вот так то лучше
	Amount uint   `json:"amount" example:"100"                              description:"Сумма перевода"                           validate:"required"` //nolint:lll,tagalign // вот так то лучше
}

// @Summary     Перевод средств с одного кошелька на другой
// @Tags  	    Wallet
// @Param walletId path string true "ID кошелька"
// @Param input body transactionRequest true "Запрос перевода средств"
// @Success     200 "Перевод успешно проведен"
// @Failure     400 "Ошибка в пользовательском запросе"
// @Failure     404 "Исходящий кошелек не найден"
// @Failure     500 "Ошибка перевода"
// @Failure     504 "Время ожидания вышло"
// @Router      /wallet/{walletId}/send [post].
func (r *walletRoutes) sendFunds(c *gin.Context) {
	var transactionRequest transactionRequest

	if err := c.BindJSON(&transactionRequest); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	walletID := c.Param("walletId")

	err := r.w.SendFunds(c.Request.Context(), walletID, transactionRequest.To, transactionRequest.Amount)
	if err != nil {
		if errors.Is(err, entity.ErrSenderIsReceiver) ||
			errors.Is(err, entity.ErrWrongAmount) ||
			errors.Is(err, entity.ErrEmptyWallet) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if errors.Is(err, entity.ErrWalletNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if errors.Is(err, entity.ErrTimeout) {
			c.AbortWithStatus(http.StatusGatewayTimeout)
			return
		}

		r.l.Error("http - v1 - sendFunds", logger.Err(err))
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	c.Status(http.StatusOK)
}

// @Summary     Получение историй входящих и исходящих транзакций
// @Description Возвращает историю транзакций по указанному кошельку.
// @Tags  	    Wallet
// @Param walletId path string true "ID кошелька"
// @Success     200 {object} []entity.Transaction "История транзакций получена"
// @Failure     404 "Указанный кошелек не найден"
// @Failure     500 "Не удалось выполнить запрос"
// @Failure     504 "Время ожидания вышло"
// @Router      /wallet/{walletId}/history [get].
func (r *walletRoutes) GetWalletHistoryByID(c *gin.Context) {
	transactions, err := r.w.GetWalletHistoryByID(c.Request.Context(), c.Param("walletId"))
	if err != nil {
		if errors.Is(err, entity.ErrTimeout) {
			c.AbortWithStatus(http.StatusGatewayTimeout)
			return
		}

		if errors.Is(err, entity.ErrWalletNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		r.l.Error("http - v1 - GetWalletHistoryByID", logger.Err(err))
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	c.JSON(http.StatusOK, transactions)
}

// @Summary     Получение текущего состояния кошелька
// @Tags  	    Wallet
// @Param walletId path string true "ID кошелька"
// @Success     200 {object} entity.Wallet "OK"
// @Failure     404 "Указанный кошелек не найден"
// @Failure     500 "Не удалось выполнить запрос"
// @Failure     504 "Время ожидания вышло"
// @Router      /wallet/{walletId} [get].
func (r *walletRoutes) GetWalletByID(c *gin.Context) {
	wallet, err := r.w.GetWalletByID(c.Request.Context(), c.Param("walletId"))
	if err != nil {
		if errors.Is(err, entity.ErrTimeout) {
			c.AbortWithStatus(http.StatusGatewayTimeout)
			return
		}

		if errors.Is(err, entity.ErrWalletNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		r.l.Error("http - v1 - GetWalletByID", logger.Err(err))
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	c.JSON(http.StatusOK, wallet)
}
