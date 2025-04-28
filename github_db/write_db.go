package githubdb

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/google/go-github/github"
	"github.com/koulipfs/auth"
	"github.com/koulipfs/constants"
	helperfunc "github.com/koulipfs/helper_func"
	"github.com/koulipfs/model"
)

// appends a new transaction to the current pool creating a new pool if full.
func WriteToGitHub(newEntry model.Transaction) error {
	ctx := context.Background()
	client, err := auth.GitAuth()
	if err != nil {
		return err
	}

	//   poolNumber and cap.
	poolNumber, ok, err := helperfunc.PoolHasCapacity()
	if err != nil {
		return err
	}
	if !ok {
		poolNumber, err = CreatePoolGitDB()
		if err != nil {
			return err
		}
	}

	filePath := path.Join(
		constants.DIR_NAME,
		fmt.Sprintf("pool_%d", poolNumber),
		constants.FILE_NAME,
	)

	//   existing content (404 empty)
	fileContent, _, _, err := client.Repositories.GetContents(
		ctx, constants.REPO_OWNER, constants.REPO_NAME, filePath, nil,
	)
	var txns []model.IPFSTransaction
	var sha *string
	if err == nil && fileContent != nil {
		sha = fileContent.SHA
		raw, _ := fileContent.GetContent()
		_ = json.Unmarshal([]byte(raw), &txns)
	}

	//  new block
	var newBlock model.IPFSTransaction
	if len(txns) == 0 {

		payload := fmt.Sprintf("%d|%s|%s|%s|%f|%d|%s",
			0,
			"",
			newEntry.TxnId,
			newEntry.ToId,
			newEntry.Amount,
			newEntry.Nonce,
			newEntry.Time,
		)

		// Generate CID/hash
		CID := helperfunc.GenerateCID(payload)

		newBlock = model.IPFSTransaction{
			CID:       CID,
			Hash:      helperfunc.GenerateHash(newEntry, ""),
			Index:     0,
			PoolIndex: poolNumber,
			TxnId:     newEntry.TxnId,
			ToId:      newEntry.ToId,
			FromId:    newEntry.FromId,
			Amount:    newEntry.Amount,
			Nonce:     newEntry.Nonce,
			Time:      newEntry.Time,
		}
	} else {
		last := txns[len(txns)-1]
		payload := fmt.Sprintf("%d|%s|%s|%s|%f|%d|%s",
			last.Index+1,
			last.Hash,
			newEntry.TxnId,
			newEntry.ToId,
			newEntry.Amount,
			newEntry.Nonce,
			newEntry.Time,
		)

		// Generate CID/hash
		CID := helperfunc.GenerateCID(payload)
		newBlock = model.IPFSTransaction{
			CID:       CID,
			PrevHash:  last.Hash,
			Hash:      helperfunc.GenerateHash(newEntry, last.Hash),
			Index:     last.Index + 1,
			PoolIndex: poolNumber,
			TxnId:     newEntry.TxnId,
			ToId:      newEntry.ToId,
			FromId:    newEntry.FromId,
			Amount:    newEntry.Amount,
			Nonce:     newEntry.Nonce,
			Time:      newEntry.Time,
		}
	}
	txns = append(txns, newBlock)

	content, _ := json.MarshalIndent(txns, "", "  ")
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(fmt.Sprintf(
			"Add txn to pool_%d at %s",
			poolNumber, time.Now().Format(time.RFC3339),
		)),
		Content: content,
		Branch:  github.String("main"),
		SHA:     sha,
	}

	if sha == nil {
		_, _, err = client.Repositories.CreateFile(ctx,
			constants.REPO_OWNER, constants.REPO_NAME, filePath, opts,
		)
	} else {
		_, _, err = client.Repositories.UpdateFile(ctx,
			constants.REPO_OWNER, constants.REPO_NAME, filePath, opts,
		)
	}
	return err
}
