package fileop

import "os"

//just creates the file
func CreateFile(fileName string) error {
	_, err := os.OpenFile(fileName, os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	return nil
}
