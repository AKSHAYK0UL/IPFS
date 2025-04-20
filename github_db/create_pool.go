package githubdb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/koulipfs/auth"
	"github.com/koulipfs/constants"
	"github.com/koulipfs/model"
)

// creates a new pool_N directory and a transactions.json file.
func CreatePoolGitDB() (int, error) {
	ctx := context.Background()
	client, err := auth.GitAuth()
	if err != nil {
		return 0, err
	}

	//  all existing pools (same logic as PoolHasCapacity)
	_, dirContents, _, err := client.Repositories.GetContents(ctx,
		constants.REPO_OWNER, constants.REPO_NAME, constants.DIR_NAME, nil,
	)
	if err != nil && !strings.Contains(err.Error(), "404") {
		return 0, err
	}

	//  next pool number
	maxPool := 0
	for _, c := range dirContents {
		if c.GetType() == "dir" && strings.HasPrefix(c.GetName(), "pool_") {
			var n int
			if _, err := fmt.Sscanf(c.GetName(), "pool_%d", &n); err == nil && n > maxPool {
				maxPool = n
			}
		}
	}
	newPool := maxPool + 1

	//  Create the transactions.json file in "pools/pool_<newPool>/"
	folderPath := fmt.Sprintf("%s/pool_%d", constants.DIR_NAME, newPool)
	filePath := fmt.Sprintf("%s/%s", folderPath, constants.FILE_NAME)
	initData, _ := json.Marshal([]model.IPFSTransaction{})
	msg := fmt.Sprintf("Initialize pool_%d on %s", newPool,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(msg),
		Content: initData,
		Branch:  github.String("main"),
	}
	_, _, err = client.Repositories.CreateFile(ctx,
		constants.REPO_OWNER, constants.REPO_NAME, filePath, opts,
	)
	if err != nil {
		return 0, err
	}
	return newPool, nil
}
