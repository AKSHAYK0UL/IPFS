package githubdb

import (
	"context"
	"encoding/json"

	"github.com/google/go-github/github"
	"github.com/koulipfs/auth"
	"github.com/koulipfs/constants"
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

	var githubContentTxn []model.Transaction
	if err := json.Unmarshal([]byte(remoteContent), &githubContentTxn); err != nil {

		return err
	}

	updated := append(githubContentTxn, newEntries)

	updatedContent, err := json.MarshalIndent(updated, "", "  ")
	if err != nil {

		return err
	}

	opt := &github.RepositoryContentFileOptions{
		Message: github.String("Append new user entries from local file"),
		Content: updatedContent,
		SHA:     github.String(sha),
		Branch:  github.String("main"),
	}

	_, _, err = client.Repositories.UpdateFile(ctx, constants.REPO_OWNER, constants.REPO_NAME, constants.FILE_NAME, opt)
	return err
}
