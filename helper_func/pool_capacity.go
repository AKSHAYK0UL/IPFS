// helperfunc/pools.go
package helperfunc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/koulipfs/auth"
	"github.com/koulipfs/constants"
	"github.com/koulipfs/model"
)

func PoolHasCapacity() (int, bool, error) {
	ctx := context.Background()
	client, err := auth.GitAuth()
	if err != nil {
		return 0, false, err
	}

	//  "pools/" dir contents
	_, dirContents, _, err := client.Repositories.GetContents(
		ctx, constants.REPO_OWNER, constants.REPO_NAME, constants.DIR_NAME, nil,
	)
	if err != nil {
		// If pools doesn't exist
		if strings.Contains(err.Error(), "404") {
			return 1, true, nil
		}
		return 0, false, err
	}

	//  Find max pool no.
	maxPool := 0
	for _, entry := range dirContents {
		if entry.GetType() == "dir" && strings.HasPrefix(entry.GetName(), "pool_") {
			var n int
			if _, err := fmt.Sscanf(entry.GetName(), "pool_%d", &n); err == nil && n > maxPool {
				maxPool = n
			}
		}
	}
	if maxPool == 0 {
		return 1, true, nil
	}

	path := fmt.Sprintf("%s/pool_%d/%s", constants.DIR_NAME, maxPool, constants.FILE_NAME)
	fileContent, _, _, err := client.Repositories.GetContents(
		ctx, constants.REPO_OWNER, constants.REPO_NAME, path, nil,
	)
	if err != nil {
		// 404 == no file
		if strings.Contains(err.Error(), "404") {
			return maxPool, true, nil
		}
		return 0, false, err
	}
	if fileContent == nil {
		return maxPool, true, nil
	}

	//  Count existing transactions
	raw, err := fileContent.GetContent()
	if err != nil {
		return 0, false, err
	}
	var txns []model.IPFSTransaction
	if err := json.Unmarshal([]byte(raw), &txns); err != nil {
		return 0, false, err
	}

	return maxPool, len(txns) < constants.MAX_TXN, nil
}
