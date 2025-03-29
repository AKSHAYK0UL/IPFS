package helperfunc

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/koulipfs/model"
)

func GenerateHash(txn model.Transaction, lastBlockHash string) string {
	data := fmt.Sprintf("%s%s%s%f%d%s", txn.TxnId, txn.ToId, txn.FromId, txn.Amount, txn.Nonce, txn.Time)
	hashBytes := sha256.Sum256([]byte(data))
	hashString := hex.EncodeToString(hashBytes[:])
	return hashString
}
