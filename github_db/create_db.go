package githubdb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/go-github/github"
	"github.com/koulipfs/auth"
	"github.com/koulipfs/constants"
	"github.com/koulipfs/model"
)

func CreateGitDB() error {
	ctx := context.Background()
	client, err := auth.GitAuth()
	if err != nil {

		return err
	}

	message := fmt.Sprintf("Created DB on %s", time.Now().Format("2006-01-02 15:04:05"))
	data, err := json.Marshal([]model.IPFSTransaction{})

	if err != nil {
		return err
	}
	createOpt := &github.RepositoryContentFileOptions{
		Message: github.String(message),
		Content: data,
		Branch:  github.String("main"),
	}
	_, _, err = client.Repositories.CreateFile(ctx, constants.REPO_OWNER, constants.REPO_NAME, constants.FILE_NAME, createOpt)
	if err != nil {
		return err
	}
	return nil
}
