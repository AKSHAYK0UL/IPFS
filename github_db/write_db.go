package githubdb

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"github.com/koulipfs/auth"
	"github.com/koulipfs/constants"
	"github.com/koulipfs/model"
)

func WriteToGitHub() error {
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

	var githubContentTxn []model.Transaction
	if err := json.Unmarshal([]byte(remoteContent), &githubContentTxn); err != nil {
		return err
	}

	localFile, err := os.Open(constants.FILE_NAME)
	if err != nil {
		return err
	}
	defer localFile.Close()

	var localContent []model.Transaction
	if err := json.NewDecoder(localFile).Decode(&localContent); err != nil {
		return err
	}

	userExists := func(txnId string, txn []model.Transaction) bool {
		for _, t := range txn {
			if t.TxnId == txnId {
				return true
			}
		}
		return false
	}

	var newEntries []model.Transaction
	for _, t := range localContent {
		if !userExists(t.TxnId, githubContentTxn) {
			newEntries = append(newEntries, t)
		}
	}

	if len(newEntries) == 0 {
		fmt.Println("No new entries to append.")
		return nil
	}

	updatedUsers := append(githubContentTxn, newEntries...)

	updatedContent, err := json.MarshalIndent(updatedUsers, "", "  ")
	if err != nil {
		return err
	}

	opt := &github.RepositoryContentFileOptions{
		Message: github.String("Append new user entries from local file"),
		Content: updatedContent,
		SHA:     github.String(sha),
		Branch:  github.String("main"),
	}

	_, _, err = client.Repositories.CreateFile(ctx, constants.REPO_OWNER, constants.REPO_NAME, constants.FILE_NAME, opt)
	return err
}
