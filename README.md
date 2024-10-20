# Locator

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
instance, err := servicelocator.GetMyService
if err != nil {
    // handle error
}
```
