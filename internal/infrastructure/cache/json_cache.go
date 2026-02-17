package cache

import (
	"encoding/json"
	"os"
	"sync"
)

type FileCache struct {
	mu sync.RWMutex
	filePath string
	data map[string][]string
}

func NewFileCache(path string) *FileCache {
	cache := &FileCache{
		filePath: path,
		data: make(map[string][]string),
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

func (c *FileCache) Set(id string, skills []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[id] = skills
	c.save()
}

func (c *FileCache) load() {
    file, err := os.ReadFile(c.filePath)
    if err == nil {
        json.Unmarshal(file, &c.data)
    }
}

func (c *FileCache) save() {
	data, _ := json.MarshalIndent(c.data, ""," ")
	_ = os.WriteFile(c.filePath, data, 0644)
}

