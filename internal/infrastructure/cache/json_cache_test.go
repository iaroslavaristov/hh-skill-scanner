package cache

import (
	"os"
	"reflect"
	"testing"
)

func TestFileCache_SaveAndGet(t *testing.T) {
	tmpFile := "test_skills_cache.json"
	defer os.Remove(tmpFile)

	storage := NewFileCache(tmpFile)
	testID := "golang-developer-50"
	expectedSkills := []string{"Go", "gRPC", "PostgreSQL"}

	storage.Set(testID, expectedSkills)

	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Error("Файл кэша не был создан после вызова Set")
	}

	storage2 := NewFileCache(tmpFile)
	
	loadedSkills, found := storage2.Get(testID)
	if !found {
		t.Error("Данные не найдены в кэше после перезагрузки")
	}

	if !reflect.DeepEqual(expectedSkills, loadedSkills) {
		t.Errorf("Данные искажены. Ожидалось %v, получено %v", expectedSkills, loadedSkills)
	}
}
