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
		seen := make(map[string]struct{})
		emit := func(path string) {
			resolved, err := filepath.EvalSymlinks(path)
			if err != nil {
				resolved = path
			}
			if _, ok := seen[resolved]; ok {
				return
			}
			seen[resolved] = struct{}{}
			results <- path
		}
		for root := range roots {
			err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() && info.Name() == ".git" {
					return filepath.SkipDir
				}
				if info.Mode()&os.ModeSymlink != 0 {
					target, err := os.Stat(path) // follows symlinks
					if err != nil {
						return nil // broken symlink
					}
					if isGitRepository(path, target.IsDir()) {
						emit(path)
					}
					return nil // never descend into symlinks
				}
				if isGitRepository(path, info.IsDir()) {
					emit(path)
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
