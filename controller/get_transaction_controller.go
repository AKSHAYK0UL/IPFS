package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koulipfs/helper"
)

func GetTransactionController(ctx *gin.Context) {
	response, err := helper.GetTransaction("")
	if err != nil {

		ctx.String(http.StatusInternalServerError, err.Error())
	} else {
		ctx.JSON(http.StatusOK, response)
	}

}
