package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koulipfs/helper"
)

func SmartContractController(ctx *gin.Context) {
	response, err := helper.GetTransaction("")
	if err != nil {

		ctx.String(http.StatusInternalServerError, err.Error())
	} else {
		invalidBlocks, err := helper.SmartContract(response)
		if err != nil {
			if len(invalidBlocks) > 0 {
				ctx.JSON(http.StatusInternalServerError, invalidBlocks)
			} else {
				ctx.String(http.StatusInternalServerError, err.Error())
			}
		} else if len(invalidBlocks) == 0 {
			ctx.JSON(http.StatusOK, gin.H{"message": "chain is valid"})

		} else {
			ctx.JSON(http.StatusOK, response)
		}

	}

}
