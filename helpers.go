package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

const ShellToUse = "bash"

func ensureDir(dirName string) error {

	err := os.Mkdir(dirName, 0700)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}
func Shellout(command string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Run()

}

func RemoveAllFolderInVideo(pathVideo string) {
	dir, _ := ioutil.ReadDir(pathVideo)
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{"video", d.Name()}...))
	}
}

func CountFileInFolder(path string) int {
	files, _ := ioutil.ReadDir(path)

	return len(files)
}
func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
