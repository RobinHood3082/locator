# Locator
[![Go Reference](https://pkg.go.dev/badge/github.com/RobinHood3082/locator.svg)](https://pkg.go.dev/github.com/RobinHood3082/locator)

The Service Locator pattern is implemented in this Go package to manage service registration and retrieval. This package allows you to register services as singletons, lazy singletons, or factories, and retrieve them as needed.

## Installation

To install the package, run:

```sh
go get github.com/RobinHood3082/locator
```

## Usage
### Creating a Service Locator
To create a new instance of the `ServiceLocator`, use the `New` function:
```go
import "github.com/RobinHood3082/locator"

sl := locator.New()
```
### Registering Services
#### Registering a Singleton
To register an already created instance as a singleton:
```go
locator.RegisterSingleton(sl, myInstance)
```
#### Registering a Lazy Singleton
To register a provider function that will create a singleton instance on first access:
```go
locator.RegisterLazySingleton(sl, func() *MyService {
    return &NewMyService()
})
```
#### Registering a Factory
To register a provider function that will create a new instance each time `Get` is called:
```go
locator.RegisterFactory(sl, func() *MyService {
    return &NewMyService()
})
```
#### Retrieving Services
To retrieve an instance of the requested type:
```go
instance, err := locator.Get[MyService](sl)
if err != nil {
    // handle error
}
```
## Example
An example is given below:
```go

import (
	"fmt"

	slocator "github.com/RobinHood3082/locator"
)

type MyService struct {
	x int
}

type MyService1 struct {
	x int
	y int
}

func NewService(x int) *MyService {
	return &MyService{x: x}
}

func NewService1(x int, y int) *MyService1 {
	return &MyService1{x: x, y: y}
}

func (svc *MyService1) Sum() int {
	return svc.x + svc.y
}

func (svc *MyService1) ChangeX(x int) {
	svc.x = x
}

func main() {
	sl := slocator.New()

	slocator.RegisterLazySingleton(sl, func() *MyService {
		return NewService(10)
	})

	slocator.RegisterFactory(sl, func() *MyService1 {
		return NewService1(10, 20)
	})

	// svc is a singleton that is created lazily
	svc, err := slocator.Get[*MyService](sl)
	if err != nil {
		fmt.Printf("Error %s\n", err)
		return
	}
	fmt.Printf("%d\n", svc.x)

	svc.x = 30

	// svcAgain is the same instance as svc
	svcAgain, err := slocator.Get[*MyService](sl)
	if err != nil {
		fmt.Printf("Error %s\n", err)
		return
	}
	fmt.Printf("%d\n", svcAgain.x) // 30

	// svc1 and svc2 are different instances
	svc1, err := slocator.Get[*MyService1](sl)
	if err != nil {
		fmt.Printf("Error %s\n", err)
		return
	}

	svc2, err := slocator.Get[*MyService1](sl)
	if err != nil {
		fmt.Printf("Error %s\n", err)
		return
	}

	// svc1 and svc2 have different values
	svc1.ChangeX(30)
	fmt.Printf("%d\n", svc1.Sum()) // 50
	fmt.Printf("%d\n", svc2.Sum()) // 30

	// Register a singleton instance
	myServiceSingleton := NewService(50)
	slocator.RegisterSingleton(sl, myServiceSingleton)

	// Retrieve the singleton instance
	svcSingleton, err := slocator.Get[*MyService](sl)
	if err != nil {
		fmt.Printf("Error %s\n", err)
		return
	}
	fmt.Printf("%d\n", svcSingleton.x) // 50

	// Modify the singleton instance
	svcSingleton.x = 100

	// Retrieve the singleton instance again
	svcSingletonAgain, err := slocator.Get[*MyService](sl)
	if err != nil {
		fmt.Printf("Error %s\n", err)
		return
	}
	fmt.Printf("%d\n", svcSingletonAgain.x) // 100
}
```
