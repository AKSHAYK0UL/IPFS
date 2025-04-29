package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/go-github/github"
	"github.com/koulipfs/auth"
	"github.com/koulipfs/constants"
	helperfunc "github.com/koulipfs/helper_func"
	"github.com/koulipfs/model"
)

func SmartContract(txns []model.IPFSTransaction) ([]model.IPFSTransaction, error) {
	if len(txns) == 0 {
		return nil, errors.New("transaction list is empty")
	}

	// 1. Group transactions by pool and validate
	pools := make(map[int][]model.IPFSTransaction)
	for _, blk := range txns {
		pools[blk.PoolIndex] = append(pools[blk.PoolIndex], blk)
	}

	var invalid []model.IPFSTransaction
	for _, group := range pools {
		prevHash := ""
		for i, blk := range group {
			computed := helperfunc.GenerateHash(model.Transaction{
				TxnId:  blk.TxnId,
				ToId:   blk.ToId,
				FromId: blk.FromId,
				Amount: blk.Amount,
				Nonce:  blk.Nonce,
				Time:   blk.Time,
			}, prevHash)

			if computed != blk.Hash || (i > 0 && blk.PrevHash != prevHash) {
				invalid = append(invalid, blk)
				break // Stop validating this pool after first invalid block
			}
			prevHash = blk.Hash
		}
	}

	ctx := context.Background()
	client, err := auth.GitAuth()
	if err != nil {
		return nil, err
	}

	if len(invalid) > 0 {
		// Track affected pools
		invalidPools := make(map[int]bool)
		for _, inv := range invalid {
			invalidPools[inv.PoolIndex] = true
		}

		// Record invalid transactions
		invalidData, _ := json.MarshalIndent(invalid, "", "  ")
		invalidPath := fmt.Sprintf("%s/invalid_%d.json", constants.DIR_NAME, time.Now().Unix())
		_, _, err = client.Repositories.CreateFile(ctx,
			constants.REPO_OWNER, constants.REPO_NAME, invalidPath,
			&github.RepositoryContentFileOptions{
				Message: github.String(fmt.Sprintf("Invalid transactions at %s", time.Now().Format(time.RFC3339))),
				Content: invalidData,
				Branch:  github.String("main"),
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to log invalid txns: %w", err)
		}

		// Load snapshot data
		snapshotData, err := helperfunc.LoadLatestSnapshot(ctx, client)
		if err != nil {
			return nil, fmt.Errorf("failed to load snapshot: %w", err)
		}

		// Group snapshot data by pool
		snapshotPools := make(map[int][]model.IPFSTransaction)
		for _, blk := range snapshotData {
			snapshotPools[blk.PoolIndex] = append(snapshotPools[blk.PoolIndex], blk)
		}

		// Process each invalid pool
		for poolIdx := range invalidPools {
			poolPath := fmt.Sprintf("%s/pool_%d/%s", constants.DIR_NAME, poolIdx, constants.FILE_NAME)

			// Delete existing file (optional, can be skipped if updating with SHA)
			// Not mandatory since we'll update or create based on SHA check

			// Recreate or update from snapshot if data exists
			if poolData, exists := snapshotPools[poolIdx]; exists {
				content, _ := json.MarshalIndent(poolData, "", "  ")

				// Check if file exists to get its SHA
				fileContent, _, resp, _ := client.Repositories.GetContents(ctx,
					constants.REPO_OWNER, constants.REPO_NAME, poolPath, &github.RepositoryContentGetOptions{Ref: "main"},
				)

				if resp != nil && resp.StatusCode == 200 && fileContent != nil && fileContent.SHA != nil {
					// File exists — update it using SHA
					_, _, err = client.Repositories.UpdateFile(ctx,
						constants.REPO_OWNER, constants.REPO_NAME, poolPath,
						&github.RepositoryContentFileOptions{
							Message: github.String(fmt.Sprintf("Restore pool_%d from snapshot", poolIdx)),
							Content: content,
							SHA:     fileContent.SHA,
							Branch:  github.String("main"),
						},
					)
					if err != nil {
						return nil, fmt.Errorf("failed to update pool_%d: %w", poolIdx, err)
					}
				} else {
					// File doesn't exist — create it
					_, _, err = client.Repositories.CreateFile(ctx,
						constants.REPO_OWNER, constants.REPO_NAME, poolPath,
						&github.RepositoryContentFileOptions{
							Message: github.String(fmt.Sprintf("Restore pool_%d from snapshot", poolIdx)),
							Content: content,
							Branch:  github.String("main"),
						},
					)
					if err != nil {
						return nil, fmt.Errorf("failed to create pool_%d: %w", poolIdx, err)
					}
				}
			}
		}

		return nil, fmt.Errorf("invalid transactions found and fixed — rolled back %d pools to snapshot", len(invalidPools))
	}

	// Create new snapshot if validation passed
	if err := helperfunc.CreateSignedSnapshot(ctx, client, txns); err != nil {
		return nil, fmt.Errorf("snapshot creation failed: %w", err)
	}

	return nil, nil
}
