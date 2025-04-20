package githubdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/koulipfs/auth"
	"github.com/koulipfs/constants"
	"github.com/koulipfs/model"
)

// GetTransaction reads all transactions from every pool folder under "pools/"
// Optionally filters by txn ID if id != "".
func GetTransaction(id string) ([]model.IPFSTransaction, error) {
	ctx := context.Background()
	client, err := auth.GitAuth()
	if err != nil {
		return nil, err
	}

	// "pools/" dir
	_, dirContents, _, err := client.Repositories.GetContents(
		ctx, constants.REPO_OWNER, constants.REPO_NAME, constants.DIR_NAME, nil,
	)
	if err != nil {
		return nil, err
	}
	if dirContents == nil {
		return nil, errors.New("no pools directory found")
	}

	var allTxns []model.IPFSTransaction
	for _, entry := range dirContents {
		if entry.GetType() != "dir" || !strings.HasPrefix(entry.GetName(), "pool_") {
			continue
		}

		filePath := fmt.Sprintf(
			"%s/%s/%s",
			constants.DIR_NAME,
			entry.GetName(),
			constants.FILE_NAME,
		)

		fileContent, _, _, err := client.Repositories.GetContents(
			ctx, constants.REPO_OWNER, constants.REPO_NAME, filePath, nil,
		)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				continue
			}
			return nil, err
		}
		if fileContent == nil {
			continue
		}

		raw, err := fileContent.GetContent()
		if err != nil {
			return nil, err
		}

		var txns []model.IPFSTransaction
		if err := json.Unmarshal([]byte(raw), &txns); err != nil {
			return nil, err
		}
		allTxns = append(allTxns, txns...)
	}

	//  If an ID is provided, filter for that single transaction
	if id != "" {
		for _, t := range allTxns {
			if t.TxnId == id {
				return []model.IPFSTransaction{t}, nil
			}
		}
		return nil, errors.New("no transaction found for ID: " + id)
	}

	return allTxns, nil
}
