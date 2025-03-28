package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koulipfs/helper"
)

func GetByIdTransactionController(ctx *gin.Context) {
	id := ctx.Param("id")
	response, err := helper.GetTransaction(id)
	if err != nil {

		ctx.String(http.StatusInternalServerError, err.Error())
	} else {
		ctx.JSON(http.StatusOK, response)
	}

}
