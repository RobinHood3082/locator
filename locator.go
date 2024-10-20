package locator

import (
	"fmt"
	"reflect"
	"sync"
)

// Provider is a function type that creates instances of services
type Provider[T any] func() T

// ServiceLocator manages service registration and retrieval
type ServiceLocator struct {
	mu        sync.RWMutex
	instances map[any]any
	providers map[any]any
}

// New creates a new ServiceLocator instance
func New() *ServiceLocator {
	return &ServiceLocator{
		instances: make(map[any]any),
		providers: make(map[any]any),
	}
}

// RegisterSingleton registers an already created instance as a singleton
func RegisterSingleton[T any](sl *ServiceLocator, instance T) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.instances[getTypeKey[T]()] = instance
}

// RegisterLazySingleton registers a provider function that will be used to create
// a singleton instance on first access
func RegisterLazySingleton[T any](sl *ServiceLocator, provider Provider[T]) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.providers[getTypeKey[T]()] = &lazySingleton[T]{
		provider: provider,
	}
}

// RegisterFactory registers a provider function that will create a new instance
// each time Get is called
func RegisterFactory[T any](sl *ServiceLocator, provider Provider[T]) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.providers[getTypeKey[T]()] = provider
}

// Get retrieves an instance of the requested type
func Get[T any](sl *ServiceLocator) (T, error) {
	sl.mu.RLock()
	typeKey := getTypeKey[T]()

	if instance, exists := sl.instances[typeKey]; exists {
		sl.mu.RUnlock()
		return instance.(T), nil
	}

	if provider, exists := sl.providers[typeKey]; exists {
		sl.mu.RUnlock()

		switch p := provider.(type) {
		case *lazySingleton[T]:
			return p.getInstance(sl)
		case Provider[T]:
			return p(), nil
		}
	} else {
		sl.mu.RUnlock()
	}

	var zero T
	return zero, fmt.Errorf("no provider registered for type %T", zero)
}

// lazySingleton wraps a provider function and ensures only one instance is created
type lazySingleton[T any] struct {
	once     sync.Once
	instance T
	provider Provider[T]
}

// getInstance returns the singleton instance, creating it if necessary
func (ls *lazySingleton[T]) getInstance(sl *ServiceLocator) (T, error) {
	if ls.provider == nil {
		var zero T
		return zero, fmt.Errorf("no provider registered for type %T", ls.instance)
	}

	ls.once.Do(func() {
		ls.instance = ls.provider()
		sl.mu.Lock()
		sl.instances[getTypeKey[T]()] = ls.instance
		sl.mu.Unlock()
	})
	return ls.instance, nil
}

// getTypeKey returns a unique key for type T
func getTypeKey[T any]() any {
	var zero T
	return reflect.TypeOf(zero)
}
