package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/xiaowuzai/simplebank/db/sqlc"
	"github.com/xiaowuzai/simplebank/token"
)

type transferRequest struct {
	FromAccountID int64  `json:"fromAccountID" binding:"required,min=1"`
	ToAccountID   int64  `json:"toAccountID" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`       // 大于零
	Currency      string `json:"currency" binding:"required,currency"` //要与账户类型匹配
}

func (s *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fromAccount, valid := s.valiadAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	// 权限验证
	authPayload := ctx.MustGet(ctxPayloadKey).(*token.Payload)
	if authPayload.Username != fromAccount.Owner {
		err := fmt.Errorf("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = s.valiadAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	// authPayload := ctx.MustGet(ctxPayloadKey).(*token.Payload)
	// if authPayload.Username != "" {

	// }

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := s.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (s *Server) valiadAccount(ctx *gin.Context, accountId int64, currency string) (db.Account, bool) {
	account, err := s.store.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}
	return account, true
}
