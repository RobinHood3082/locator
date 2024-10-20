package locator_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/RobinHood3082/locator"
)

type TestService struct {
	Name string
}

type AnotherTestService struct {
	ID int
}

// Test RegisterSingleton and Get
func TestSingleton(t *testing.T) {
	sl := locator.New()

	singletonInstance := &TestService{Name: "Singleton"}
	locator.RegisterSingleton(sl, singletonInstance)

	retrievedSingleton, err := locator.Get[*TestService](sl)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if retrievedSingleton != singletonInstance {
		t.Fatalf("expected %v, got %v", singletonInstance, retrievedSingleton)
	}
}

// Test RegisterLazySingleton and Get
func TestLazySingleton(t *testing.T) {
	sl := locator.New()

	lazySingletonProvider := func() *TestService {
		return &TestService{Name: "LazySingleton"}
	}
	locator.RegisterLazySingleton(sl, lazySingletonProvider)

	retrievedLazySingleton, err := locator.Get[*TestService](sl)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if retrievedLazySingleton.Name != "LazySingleton" {
		t.Fatalf("expected LazySingleton, got %v", retrievedLazySingleton.Name)
	}
}

// Test RegisterFactory and Get
func TestFactory(t *testing.T) {
	sl := locator.New()

	factoryProvider := func() *TestService {
		return &TestService{Name: "Factory"}
	}
	locator.RegisterFactory(sl, factoryProvider)

	retrievedFactory1, err := locator.Get[*TestService](sl)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if retrievedFactory1.Name != "Factory" {
		t.Fatalf("expected Factory, got %v", retrievedFactory1.Name)
	}

	retrievedFactory2, err := locator.Get[*TestService](sl)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if retrievedFactory2.Name != "Factory" {
		t.Fatalf("expected Factory, got %v", retrievedFactory2.Name)
	}

	if retrievedFactory1 == retrievedFactory2 {
		t.Fatalf("expected different instances, got the same")
	}
}

// Test Get with unregistered type
func TestUnregisteredType(t *testing.T) {
	sl := locator.New()

	type UnregisteredService struct{}
	_, err := locator.Get[*UnregisteredService](sl)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	expectedError := "no provider registered for type *locator_test.UnregisteredService"
	if err.Error() != expectedError {
		t.Fatalf("expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

// Test LazySingleton concurrency
func TestLazySingletonConcurrency(t *testing.T) {
	sl := locator.New()

	var callCount int
	var mu sync.Mutex

	lazySingletonProvider := func() *TestService {
		mu.Lock()
		defer mu.Unlock()
		callCount++
		return &TestService{Name: "LazySingleton"}
	}
	locator.RegisterLazySingleton(sl, lazySingletonProvider)

	var wg sync.WaitGroup
	const goroutines = 100
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			_, err := locator.Get[*TestService](sl)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		}()
	}

	wg.Wait()

	if callCount != 1 {
		t.Fatalf("expected provider to be called once, got %d", callCount)
	}
}

// Test multiple registrations and retrievals for different types
func TestMultipleRegistrations(t *testing.T) {
	sl := locator.New()

	locator.RegisterSingleton(sl, &TestService{Name: "Service1"})
	locator.RegisterSingleton(sl, &AnotherTestService{ID: 1})

	service1, err := locator.Get[*TestService](sl)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if service1.Name != "Service1" {
		t.Fatalf("expected Service1, got %v", service1.Name)
	}

	service2, err := locator.Get[*AnotherTestService](sl)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if service2.ID != 1 {
		t.Fatalf("expected ID 1, got %v", service2.ID)
	}
}

// Test overwriting a registered service
func TestOverwriteRegistration(t *testing.T) {
	sl := locator.New()

	locator.RegisterSingleton(sl, &TestService{Name: "Original"})
	locator.RegisterSingleton(sl, &TestService{Name: "Overwritten"})

	service, err := locator.Get[*TestService](sl)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if service.Name != "Overwritten" {
		t.Fatalf("expected Overwritten, got %v", service.Name)
	}
}

// Test retrieving a lazy singleton multiple times
func TestLazySingletonMultipleRetrievals(t *testing.T) {
	sl := locator.New()

	var callCount int
	locator.RegisterLazySingleton(sl, func() *TestService {
		callCount++
		return &TestService{Name: fmt.Sprintf("Instance%d", callCount)}
	})

	for i := 0; i < 3; i++ {
		service, err := locator.Get[*TestService](sl)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if service.Name != "Instance1" {
			t.Fatalf("expected Instance1, got %v", service.Name)
		}
	}

	if callCount != 1 {
		t.Fatalf("expected provider to be called once, got %d", callCount)
	}
}

// Test registering and retrieving value types
func TestValueTypes(t *testing.T) {
	sl := locator.New()

	locator.RegisterSingleton(sl, 42)
	locator.RegisterSingleton(sl, "hello")

	intValue, err := locator.Get[int](sl)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if intValue != 42 {
		t.Fatalf("expected 42, got %v", intValue)
	}

	stringValue, err := locator.Get[string](sl)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if stringValue != "hello" {
		t.Fatalf("expected hello, got %v", stringValue)
	}
}

// Test concurrent access for all registration types
func TestConcurrentAccess(t *testing.T) {
	sl := locator.New()

	locator.RegisterSingleton(sl, &TestService{Name: "Singleton"})
	locator.RegisterLazySingleton(sl, func() *AnotherTestService {
		return &AnotherTestService{ID: 1}
	})
	var factoryCounter int
	var mu sync.Mutex

	locator.RegisterFactory(sl, func() int {
		mu.Lock()
		defer mu.Unlock()
		factoryCounter++
		return factoryCounter
	})

	var wg sync.WaitGroup
	const goroutines = 100
	wg.Add(goroutines * 3)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			_, err := locator.Get[*TestService](sl)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		}()
		go func() {
			defer wg.Done()
			_, err := locator.Get[*AnotherTestService](sl)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		}()
		go func() {
			defer wg.Done()
			_, err := locator.Get[int](sl)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		}()
	}

	wg.Wait()

	if factoryCounter != goroutines {
		t.Fatalf("expected factory to be called %d times, got %d", goroutines, factoryCounter)
	}
}

// Test registering a nil provider
func TestNilProvider(t *testing.T) {
	sl := locator.New()

	locator.RegisterLazySingleton[*TestService](sl, nil)

	_, err := locator.Get[*TestService](sl)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	expectedError := "no provider registered for type *locator_test.TestService"
	if err.Error() != expectedError {
		t.Fatalf("expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

// Test type safety
func TestTypeSafety(t *testing.T) {
	sl := locator.New()

	locator.RegisterSingleton(sl, &TestService{Name: "Service"})

	_, err := locator.Get[*AnotherTestService](sl)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	expectedError := "no provider registered for type *locator_test.AnotherTestService"
	if err.Error() != expectedError {
		t.Fatalf("expected error message '%s', got '%s'", expectedError, err.Error())
	}
}
