package helperfunc

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
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

// func LoadLatestSnapshot(ctx context.Context, client *github.Client) ([]model.IPFSTransaction, error) {
// 	// get directory listing
// 	_, dirContents, _, err := client.Repositories.GetContents(ctx,
// 		constants.REPO_OWNER, constants.REPO_NAME,
// 		snapshotDir, nil,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// find max version file
// 	maxVer := 0
// 	var latestFile *github.RepositoryContent
// 	for _, entry := range dirContents {
// 		if entry.GetType() == "file" && strings.HasSuffix(entry.GetName(), ".json") {
// 			var v int
// 			if _, err := fmt.Sscanf(entry.GetName(), "v%d.json", &v); err == nil && v > maxVer {
// 				maxVer = v
// 				latestFile = entry
// 			}
// 		}
// 	}
// 	if latestFile == nil {
// 		return nil, fmt.Errorf("no snapshots found")
// 	}
// 	// load content
// 	raw, _, _, err := client.Repositories.GetContents(ctx,
// 		constants.REPO_OWNER, constants.REPO_NAME,
// 		snapshotDir+"/"+latestFile.GetName(), nil,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	data, err := raw.GetContent()
// 	if err != nil {
// 		return nil, err
// 	}
// 	var txns []model.IPFSTransaction
// 	if err := json.Unmarshal([]byte(data), &txns); err != nil {
// 		return nil, err
// 	}
// 	return txns, nil
// }

//==============================================================================

func LoadLatestSnapshot(ctx context.Context, client *github.Client) ([]model.IPFSTransaction, error) {
	// list snapshot directory
	_, dirContents, _, err := client.Repositories.GetContents(ctx,
		constants.REPO_OWNER, constants.REPO_NAME,
		snapshotDir, nil,
	)
	if err != nil {
		return nil, fmt.Errorf("listing snapshots: %w", err)
	}

	// collect version numbers
	versions := []int{}
	for _, entry := range dirContents {
		if entry.GetType() == "file" && strings.HasSuffix(entry.GetName(), ".json") {
			var v int
			if _, err := fmt.Sscanf(entry.GetName(), "v%d.json", &v); err == nil {
				versions = append(versions, v)
			}
		}
	}
	if len(versions) == 0 {
		return nil, fmt.Errorf("no snapshots found")
	}
	sort.Sort(sort.Reverse(sort.IntSlice(versions))) // descending

	var lastErr error
	secret := os.Getenv("SNAPSHOT_HMAC_SECRET")
	if secret == "" {
		return nil, errors.New("SNAPSHOT_HMAC_SECRET not set")
	}

	// try each version from latest down
	for _, ver := range versions {
		jsonPath := fmt.Sprintf("%s/v%d.json", snapshotDir, ver)
		sigPath := fmt.Sprintf("%s/v%d.sig", snapshotDir, ver)

		// load JSON
		rawJSON, _, _, err := client.Repositories.GetContents(ctx,
			constants.REPO_OWNER, constants.REPO_NAME,
			jsonPath, nil,
		)
		if err != nil {
			lastErr = fmt.Errorf("fetching %s: %w", jsonPath, err)
			continue
		}
		jsonBytes, err := rawJSON.GetContent()
		if err != nil {
			lastErr = fmt.Errorf("reading content %s: %w", jsonPath, err)
			continue
		}

		// load signature
		rawSig, _, _, err := client.Repositories.GetContents(ctx,
			constants.REPO_OWNER, constants.REPO_NAME,
			sigPath, nil,
		)
		if err != nil {
			lastErr = fmt.Errorf("fetching %s: %w", sigPath, err)
			continue
		}
		sigHex, err := rawSig.GetContent()
		if err != nil {
			lastErr = fmt.Errorf("reading signature %s: %w", sigPath, err)
			continue
		}

		// verify HMAC
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(jsonBytes))
		expected := hex.EncodeToString(mac.Sum(nil))
		if !hmac.Equal([]byte(expected), []byte(sigHex)) {
			lastErr = fmt.Errorf("signature mismatch for v%d", ver)
			continue
		}

		// unmarshal and return
		var txns []model.IPFSTransaction
		if err := json.Unmarshal([]byte(jsonBytes), &txns); err != nil {
			lastErr = fmt.Errorf("unmarshal v%d: %w", ver, err)
			continue
		}

		return txns, nil
	}

	return nil, fmt.Errorf("all snapshots invalid, last error: %v", lastErr)
}
