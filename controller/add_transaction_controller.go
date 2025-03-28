package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koulipfs/helper"
	"github.com/koulipfs/model"
)

func AddTransactionController(ctx *gin.Context) {
	txn := new(model.Transaction)
	if err := ctx.ShouldBindBodyWithJSON(txn); err != nil {

		ctx.String(http.StatusBadRequest, err.Error())
	} else {
		err := helper.AddTransaction(*txn)
		if err != nil {

			ctx.String(http.StatusInternalServerError, err.Error())
		} else {
			ctx.String(http.StatusOK, "")
		}
	}

}
