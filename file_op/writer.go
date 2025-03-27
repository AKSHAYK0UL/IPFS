package fileop

import (
	"encoding/json"
	"os"

	"github.com/koulipfs/constants"
	"github.com/koulipfs/model"
)

func WriteFile(Txn model.Transaction) error {

	fileData, err := ReadFile()
	if err != nil {

		return err
	}

	fileData = append(fileData, Txn)

	byteSlice, err := json.MarshalIndent(fileData, "", " ")

	if err != nil {
		return err
	}
	if err := os.WriteFile(constants.FILE_NAME, byteSlice, 0664); err != nil {
		return err
	}

	return nil
}
