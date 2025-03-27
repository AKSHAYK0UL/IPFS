package helper

import (
	githubdb "github.com/koulipfs/github_db"
	"github.com/koulipfs/model"
)

func GetTransaction() ([]model.Transaction, error) {

	txns, err := githubdb.GetTransaction()
	if err != nil {
		return nil, err
	}
	return txns, nil
}
