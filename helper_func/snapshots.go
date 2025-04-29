package helperfunc

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/koulipfs/constants"
	"github.com/koulipfs/model"
)

const snapshotDir = constants.DIR_SNAPSHOT

func CreateSignedSnapshot(ctx context.Context, client *github.Client, txns []model.IPFSTransaction) error {
	// 1. Determine next version
	_, dirContents, _, err := client.Repositories.GetContents(ctx,
		constants.REPO_OWNER, constants.REPO_NAME,
		snapshotDir, nil,
	)
	if err != nil && !strings.Contains(err.Error(), "404") {
		return err
	}
	// find max version among files named like "v1.json", "v2.json"
	maxVer := 0
	for _, entry := range dirContents {
		if entry.GetType() == "file" && strings.HasSuffix(entry.GetName(), ".json") {
			var v int
			if _, err := fmt.Sscanf(entry.GetName(), "v%d.json", &v); err == nil && v > maxVer {
				maxVer = v
			}
		}
	}
	newVer := maxVer + 1

	// 2. Marshal snapshot
	content, err := json.MarshalIndent(txns, "", "  ")
	if err != nil {
		return err
	}

	// 3. Compute signature (HMAC-SHA256 here)
	mac := hmac.New(sha256.New, []byte(os.Getenv("SNAPSHOT_HMAC_SECRET")))
	mac.Write(content)
	sig := hex.EncodeToString(mac.Sum(nil))

	// 4. Commit snapshot JSON
	versionedPath := fmt.Sprintf("%s/v%d.json", snapshotDir, newVer)
	msg1 := fmt.Sprintf("Snapshot v%d at %s", newVer, time.Now().Format(time.RFC3339))
	opts1 := &github.RepositoryContentFileOptions{
		Message: github.String(msg1),
		Content: content,
		Branch:  github.String("main"),
	}
	if _, _, err := client.Repositories.CreateFile(ctx, constants.REPO_OWNER, constants.REPO_NAME, versionedPath, opts1); err != nil {
		return err
	}

	// 5. Commit signature alongside
	sigPath := fmt.Sprintf("%s/v%d.sig", snapshotDir, newVer)
	msg2 := fmt.Sprintf("Signature for snapshot v%d", newVer)
	opts2 := &github.RepositoryContentFileOptions{
		Message: github.String(msg2),
		Content: []byte(sig),
		Branch:  github.String("main"),
	}
	_, _, err = client.Repositories.CreateFile(ctx, constants.REPO_OWNER, constants.REPO_NAME, sigPath, opts2)
	return err
}

func LoadLatestSnapshot(ctx context.Context, client *github.Client) ([]model.IPFSTransaction, error) {
	// get directory listing
	_, dirContents, _, err := client.Repositories.GetContents(ctx,
		constants.REPO_OWNER, constants.REPO_NAME,
		snapshotDir, nil,
	)
	if err != nil {
		return nil, err
	}
	// find max version file
	maxVer := 0
	var latestFile *github.RepositoryContent
	for _, entry := range dirContents {
		if entry.GetType() == "file" && strings.HasSuffix(entry.GetName(), ".json") {
			var v int
			if _, err := fmt.Sscanf(entry.GetName(), "v%d.json", &v); err == nil && v > maxVer {
				maxVer = v
				latestFile = entry
			}
		}
	}
	if latestFile == nil {
		return nil, fmt.Errorf("no snapshots found")
	}
	// load content
	raw, _, _, err := client.Repositories.GetContents(ctx,
		constants.REPO_OWNER, constants.REPO_NAME,
		snapshotDir+"/"+latestFile.GetName(), nil,
	)
	if err != nil {
		return nil, err
	}
	data, err := raw.GetContent()
	if err != nil {
		return nil, err
	}
	var txns []model.IPFSTransaction
	if err := json.Unmarshal([]byte(data), &txns); err != nil {
		return nil, err
	}
	return txns, nil
}
