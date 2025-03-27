package helper

import (
	"fmt"

	fileop "github.com/koulipfs/file_op"
	githubdb "github.com/koulipfs/github_db"
	"github.com/koulipfs/model"
)

func AddTransaction(txn model.Transaction) error {
	if err := fileop.WriteFile(txn); err != nil {
		return err
	}
	if err := githubdb.WriteToGitHub(); err != nil {
		fmt.Println("ERROR IN GIT DB ", err.Error())
		return err
	}
	return nil
}
