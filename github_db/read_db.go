package githubdb

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/koulipfs/auth"
	"github.com/koulipfs/constants"
	"github.com/koulipfs/model"
)

func GetTransaction(id string) ([]model.IPFSTransaction, error) {
	ctx := context.Background()
	client, err := auth.GitAuth()
	if err != nil {

		return []model.IPFSTransaction{}, err
	}

	fileContent, _, _, err := client.Repositories.GetContents(ctx, constants.REPO_OWNER, constants.REPO_NAME, constants.FILE_NAME, nil)
	if err != nil {

		return []model.IPFSTransaction{}, err
	}
	guthubContent, err := fileContent.GetContent()
	if err != nil {

		return []model.IPFSTransaction{}, err
	}

	var txns []model.IPFSTransaction
	if err := json.Unmarshal([]byte(guthubContent), &txns); err != nil {

		return []model.IPFSTransaction{}, err
	}
	if id != "" {

		exist := false
		for _, t := range txns {
			if t.TxnId == id {
				exist = true
				return []model.IPFSTransaction{t}, nil
			}
		}
		if !exist {
			return []model.IPFSTransaction{}, errors.New("no transaction found")
		}
	}
	return txns, nil

}
