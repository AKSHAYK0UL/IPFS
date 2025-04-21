// package helper

// import (
// 	"errors"

// 	helperfunc "github.com/koulipfs/helper_func"
// 	"github.com/koulipfs/model"
// )

// func SmartContract(
// 	txns []model.IPFSTransaction,
// ) ([]model.IPFSTransaction, error) {
// 	const poolSize = 3
// 	if len(txns) == 0 {
// 		return nil, errors.New("transaction chain is empty")
// 	}

// 	numPools := (len(txns) + poolSize - 1) / poolSize
// 	invalidPool := make([]bool, numPools)

// 	for start := 0; start < len(txns); start += poolSize {
// 		end := start + poolSize
// 		if end > len(txns) {
// 			end = len(txns)
// 		}
// 		poolIdx := start / poolSize
// 		prevHash := ""

// 		for i := start; i < end; i++ {
// 			blk := txns[i]
// 			computed := helperfunc.GenerateHash(model.Transaction{
// 				TxnId:  blk.TxnId,
// 				ToId:   blk.ToId,
// 				FromId: blk.FromId,
// 				Amount: blk.Amount,
// 				Nonce:  blk.Nonce,
// 				Time:   blk.Time,
// 			}, prevHash)

// 			if computed != blk.Hash || (i != start && blk.PrevHash != prevHash) {
// 				invalidPool[poolIdx] = true
// 				break
// 			}
// 			prevHash = blk.Hash
// 		}
// 	}

// 	var invalid []model.IPFSTransaction
// 	for start := 0; start < len(txns); start += poolSize {
// 		poolIdx := start / poolSize
// 		if !invalidPool[poolIdx] {
// 			continue
// 		}
// 		end := start + poolSize
// 		if end > len(txns) {
// 			end = len(txns)
// 		}
// 		for i := start; i < end; i++ {
// 			invalid = append(invalid, txns[i])
// 		}
// 	}

// 	return invalid, nil
// }

//===========================================================================

package helper

import (
	"errors"
	"fmt"

	helperfunc "github.com/koulipfs/helper_func"
	"github.com/koulipfs/model"
)

func SmartContract(
	txns []model.IPFSTransaction,
) ([]model.IPFSTransaction, error) {
	if len(txns) == 0 {
		return nil, errors.New("transaction list is empty")
	}

	//  Group txns by PoolIndex
	pools := make(map[int][]model.IPFSTransaction)
	for _, blk := range txns {
		pools[blk.PoolIndex] = append(pools[blk.PoolIndex], blk)
	}

	//  Validate each pool
	var invalid []model.IPFSTransaction
	for _, group := range pools {
		prevHash := ""
		valid := true

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
				valid = false
				break
			}
			prevHash = blk.Hash
		}

		if !valid {
			invalid = append(invalid, group...)
		}
	}

	if len(invalid) > 0 {
		return invalid, fmt.Errorf("found %d invalid pools", len(invalid))
	}
	return nil, nil // all pools valid
}
