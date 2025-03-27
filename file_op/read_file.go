package fileop

import (
	"encoding/json"
	"os"

	"github.com/koulipfs/constants"
	"github.com/koulipfs/model"
)

func ReadFile() ([]model.Transaction, error) {
	bytes, err := os.ReadFile(constants.FILE_NAME)
	if err != nil {
		if os.IsNotExist(err) {
			return []model.Transaction{}, nil
		}
		return nil, err
	}
	if len(bytes) == 0 {
		return []model.Transaction{}, nil
	}

	var txns []model.Transaction
	if err := json.Unmarshal(bytes, &txns); err == nil {
		return txns, nil
	}

	var singleTxn model.Transaction
	if err := json.Unmarshal(bytes, &singleTxn); err != nil {
		return nil, err
	}
	return []model.Transaction{singleTxn}, nil
}
