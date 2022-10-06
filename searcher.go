package repofinder

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Search(paths []string) error {
	cache, err := readCache()
	if err != nil {
		return fmt.Errorf("failed to read cache: :%w", err)
	}
	for _, sourcePath := range paths {
		if repos, ok := cache[sourceDir(sourcePath)]; ok {
			for repo := range repos {
				if _, err := os.Lstat(string(repo)); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "removing %s from cache due to %v", repo, err)
					delete(repos, repo)
					continue
				}

				printPath(repo)
				repos[repo] = true
			}
		} else {
			cache[sourceDir(sourcePath)] = make(repo)
		}
	}

	for _, rootPath := range paths {
		err := filepath.WalkDir(rootPath, func(currPath string, d fs.DirEntry, err error) error {
			if err != nil {
				if errors.Is(err, os.ErrPermission) {
					_, _ = fmt.Fprintf(os.Stderr, "error: permission denied on %s: %v\n", currPath, err)
					return filepath.SkipDir
				}
				return err
			}
			if !d.IsDir() {
				return nil
			}
			if isHidden(currPath) || isIgnored(currPath) {
				return filepath.SkipDir
			}

			absPath, err := absPath(currPath)
			if err != nil {
				return err
			}

			if alreadyPrinted := cache[sourceDir(rootPath)][repoPath(absPath)]; alreadyPrinted {
				return filepath.SkipDir
			}
			if isRepo(currPath) {
				printPath(absPath)
				cache[sourceDir(rootPath)][repoPath(absPath)] = true
				return filepath.SkipDir
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to index %s: %w", rootPath, err)
		}
	}
	return writeCache(cache)
}

func printPath(p repoPath) {
	_, _ = fmt.Fprintf(os.Stdout, "%s\n", string(p))
}

func isIgnored(p string) bool {
	if strings.HasSuffix(p, "pkg/mod") {
		// Ignore $GOPATH/pkg/mod folder
		return true
	}

	ignoredNames := []string{
		"node_modules",
		"build",
	}
	name := path.Base(p)
	for _, ignoredName := range ignoredNames {
		if name == ignoredName {
			return true
		}
	}
	return false
}

func isHidden(p string) bool {
	return strings.HasPrefix(path.Base(p), ".")
}

func isRepo(p string) bool {
	if git, err := os.Lstat(path.Join(p, ".git")); err != nil || !git.IsDir() {
		return false
	}
	return true
}

func absPath(p string) (repoPath, error) {
	absPath, err := filepath.Abs(p)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute dir of %s: %w", p, err)
	}
	return repoPath(absPath), nil
}
