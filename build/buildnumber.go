package build

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return fmt.Sprintf("%x%x%x%x%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func getNextBuildNo(repoUrl string) (int, error) {
	buildsDir := getBuildsPath(repoUrl)
	if _, err := os.Stat(buildsDir); os.IsNotExist(err) {
		return 0, nil
	}

	files, err := ioutil.ReadDir(buildsDir)
	if err != nil {
		return 0, err
	}
	return len(files) + 1, nil
}
