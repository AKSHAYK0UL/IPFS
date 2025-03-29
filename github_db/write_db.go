package githubdb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/go-github/github"
	"github.com/koulipfs/auth"
	"github.com/koulipfs/constants"
	helperfunc "github.com/koulipfs/helper_func"
	"github.com/koulipfs/model"
)

func WriteToGitHub(newEntries model.Transaction) error {
	ctx := context.Background()
	client, err := auth.GitAuth()
	if err != nil {

		return err
	}

	githubContent, _, _, err := client.Repositories.GetContents(ctx, constants.REPO_OWNER, constants.REPO_NAME, constants.FILE_NAME, nil)
	if err != nil {

		return err
	}
	sha := githubContent.GetSHA()
	remoteContent, err := githubContent.GetContent()
	if err != nil {

		return err
	}

	var githubContentTxn []model.IPFSTransaction
	if err := json.Unmarshal([]byte(remoteContent), &githubContentTxn); err != nil {

		return err
	}

	var newTxn model.IPFSTransaction
	length := len(githubContentTxn) - 1

	if length == -1 {
		newTxn = model.IPFSTransaction{
			Hash:   helperfunc.GenerateHash(newEntries, ""),
			Index:  0,
			TxnId:  newEntries.TxnId,
			ToId:   newEntries.ToId,
			FromId: newEntries.FromId,
			Amount: newEntries.Amount,
			Nonce:  newEntries.Nonce,
			Time:   newEntries.Time,
		}
	} else {
		lastBlock := githubContentTxn[length]
		newTxn = model.IPFSTransaction{
			PrevHash: lastBlock.Hash,
			Hash:     helperfunc.GenerateHash(newEntries, lastBlock.Hash),
			Index:    lastBlock.Index + 1,
			TxnId:    newEntries.TxnId,
			ToId:     newEntries.ToId,
			FromId:   newEntries.FromId,
			Amount:   newEntries.Amount,
			Nonce:    newEntries.Nonce,
			Time:     newEntries.Time,
		}

	}

	updated := append(githubContentTxn, newTxn)

	updatedContent, err := json.MarshalIndent(updated, "", "  ")
	if err != nil {

		return err
	}
	message := fmt.Sprintf("New Transaction added on %s", time.Now().Format("2006-01-02 15:04:05"))
	opt := &github.RepositoryContentFileOptions{
		Message: github.String(message),
		Content: updatedContent,
		SHA:     github.String(sha),
		Branch:  github.String("main"),
	}

	_, _, err = client.Repositories.UpdateFile(ctx, constants.REPO_OWNER, constants.REPO_NAME, constants.FILE_NAME, opt)
	return err
}
