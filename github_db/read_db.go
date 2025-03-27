package githubdb

import (
	"context"
	"encoding/json"

	"github.com/koulipfs/auth"
	"github.com/koulipfs/constants"
	"github.com/koulipfs/model"
)

func GetTransaction() ([]model.Transaction, error) {
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
	return txns, nil

}
