package cache

import (
	"encoding/json"
	"os"
	"sync"
)

type FileCache struct {
	mu       sync.RWMutex
	filePath string
	data     map[string][]string
}

func NewFileCache(path string) *FileCache {
	cache := &FileCache{
		filePath: path,
		data:     make(map[string][]string),
	}
	cache.load()
	return cache
}

func (c *FileCache) Get(id string) ([]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	skills, ok := c.data[id]
	return skills, ok
}

func (c *FileCache) Set(id string, skills []string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[id] = skills
	return nil
}

func (c *FileCache) Save() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	data, err := json.MarshalIndent(c.data, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.filePath, data, 0644)
}

func (c *FileCache) load() {
	file, err := os.ReadFile(c.filePath)
	if err == nil {
		json.Unmarshal(file, &c.data)
	}
}
