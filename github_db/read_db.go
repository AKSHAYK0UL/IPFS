package githubdb

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/koulipfs/auth"
	"github.com/koulipfs/constants"
	"github.com/koulipfs/model"
)

func GetTransaction(id string) ([]model.Transaction, error) {
	ctx := context.Background()
	client, err := auth.GitAuth()
	if err != nil {

		return []model.Transaction{}, err
	}

	fileContent, _, _, err := client.Repositories.GetContents(ctx, constants.REPO_OWNER, constants.REPO_NAME, constants.FILE_NAME, nil)
	if err != nil {

		return []model.Transaction{}, err
	}
	guthubContent, err := fileContent.GetContent()
	if err != nil {

		return []model.Transaction{}, err
	}

	var txns []model.Transaction
	if err := json.Unmarshal([]byte(guthubContent), &txns); err != nil {

		return []model.Transaction{}, err
	}
	if id != "" {

		exist := false
		for _, t := range txns {
			if t.TxnId == id {
				exist = true
				return []model.Transaction{t}, nil
			}
		}
		if !exist {
			return []model.Transaction{}, errors.New("No Transaction found!")
		}
	}
	return txns, nil

}
