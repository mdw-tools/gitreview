package main

import (
	"bufio"
	"iter"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func collectGitRepositories(roots iter.Seq[string]) (results chan string) {
	results = make(chan string)
	go func() {
		defer close(results)
		for root := range roots {
			err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() && info.Name() == ".git" {
					return filepath.SkipDir
				}
				if isGitRepository(path, info.IsDir()) {
					results <- path
					return filepath.SkipDir
				}
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
	return results
}
func isGitRepository(path string, isDir bool) bool {
	if !isDir {
		return false
	}

	_, err := os.Stat(filepath.Join(path, ".git"))
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func execute(dir, command string) (string, error) {
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func prompt(message string) {
	log.Println(message)
	bufio.NewScanner(os.Stdin).Scan()
}
