package registry

import (
	"log"
	"sync"
)

type Registry struct {
	keys map[string]any
}

var instance *Registry
var once sync.Once

func GetRegistryInstance() Registry {
	once.Do(func() {
		instance = &Registry{keys: make(map[string]interface{})}
	})
	return *instance
}

func (registry *Registry) Register(key string, value interface{}) {
	if _, ok := registry.keys[key]; ok {
		log.Fatalf("key `%s` already registered", key)
		return
	}
	registry.keys[key] = value
}

func (registry *Registry) Inject(key string) interface{} {
	if _, ok := registry.keys[key]; !ok {
		log.Fatalf("key `%s` is not registered", key)
		return nil
	}
	return registry.keys[key]
}
