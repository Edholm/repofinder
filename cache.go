package repofinder

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
)

const relativeCachePath = ".cache/repofinder/cache.json"

type (
	repoPath string
	// printed shows whether the repo has been written to stdout.
	printed   bool
	repo      map[repoPath]printed
	sourceDir string
	// cache contains previously found repos, keyed by "source search dir".
	repoCache map[sourceDir]repo
)

func getCachePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine user home dir: %w", err)
	}

	cacheDir := path.Join(homeDir, relativeCachePath)
	return cacheDir, nil
}

func readCache() (repoCache, error) {
	cachePath, err := getCachePath()
	if err != nil {
		return nil, err
	}
	cache := make(repoCache)
	cacheFd, err := os.Open(cachePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// We have no previous cache so just return an empty one.
			return cache, nil
		}
		return nil, fmt.Errorf("failed to open cache file: %w", err)
	}

	if err := json.NewDecoder(cacheFd).Decode(&cache); err != nil {
		return nil, fmt.Errorf("failed to decode cache: %w", err)
	}
	return cache, nil
}

func falsify(rc repoCache) {
	for _, r := range rc {
		for repo := range r {
			r[repo] = false
		}
	}
}

func writeCache(rc repoCache) error {
	falsify(rc)
	bytes, err := json.Marshal(rc)
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}
	cachePath, err := getCachePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(path.Dir(cachePath), 0o755); err != nil {
		return fmt.Errorf("failed to create cache dir: %w", err)
	}

	if err := os.WriteFile(cachePath, bytes, 0o600); err != nil {
		return fmt.Errorf("failed to write cache to %s: %w", cachePath, err)
	}
	return nil
}
