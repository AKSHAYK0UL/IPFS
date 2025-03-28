package helper

import (
	"fmt"

	githubdb "github.com/koulipfs/github_db"
	"github.com/koulipfs/model"
)

func AddTransaction(txn model.Transaction) error {

	if err := githubdb.WriteToGitHub(txn); err != nil {
		fmt.Println("ERROR IN GIT DB ", err.Error())
		return err
	}
	return nil
}
