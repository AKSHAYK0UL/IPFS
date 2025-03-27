package main

import (
	"os"

	"github.com/koulipfs/constants"
	fileop "github.com/koulipfs/file_op"
	"github.com/koulipfs/route"
)

// overView
// setup the github auth
// get the data
// store it into the local file in the server
// after storing the file in the server then push the file to the github(IPFS)
// ---------------------------------------------------------------------------
// step by step
// create local file
// write to it (data from the api)
// also update
// read from it (to store it in the IPFS)

func init() {

	err := fileop.CreateFile(constants.FILE_NAME)
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	route := route.RouteTable()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"

	}
	route.Run(":" + port)

}
