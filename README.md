# Storz

Storz is an Object Store framework built in golang.

It features a simple object modeling language used to [generate ](https://github.com/wazofski/storz/tree/main/mgen)
the golang object class meta.

The generated code contains object and structure classes used to interact 
with the Store interface that most storz modules expose.


## Purpose
Storz Modules can be used to store and retrieve objects from a database.

Modules can be composed together to chain Store functionality into more complex logical modules.

Such combined modules can provide ways to handle object changes and routing 
objects to different Stores of any complexity.

The interfaces are exposed from the `store` package to allow custom store implementations if needed.

Multi-client, multi-server complexity of any desired depth can be achieved with ease basically for free:

- Frontend service
  - Cached store based on..
  - ..Rest client connected to the backend service
- Backend service
  - Rest server exposing..
  - ..React store running validations and other logic on top of ..
  - ..Router store which routes 
    - Coke objects to a cached store based on..
    - ..SQL store (in another network)
    - Pepsi objects to a local MongoDB store
    - RootBeer objects to a cached..
      - ..REST client store connected to another service
- Second backend service
  - Rest server exposing..
  - ..React store running validations and other logic on top of ..
    - ..local MongoDB store


## Getting started

**Install the storz module**
```
go get github.com/wazofski/storz
go install github.com/wazofski/storz
```

**Initialize your storz project**
```
storz init [project]
```

This will create the [project] directory containing your go module, a sample model and the main.go source file.

**Generate the class meta**
In your [project] directory, run
(re-run when making changes to your model YAML files)

```
go generate
```

**Build and run your code**
```
go build
./[project]
```

## Modules
Most Storz modules are implementations of the [Store](https://github.com/wazofski/storz/tree/main/store) interface.
This allows mixing and matching various modules together into 
more complex systems.

### Persistence Layer
Persistence modules are independent, meaning they do not need  another Store to operate.
- [Memory](https://github.com/wazofski/storz/tree/main/memory) store - simple in-memory store useful for temporary storage cases
- [Mongo](https://github.com/wazofski/storz/tree/main/mongo) store - uses an existing Mongo DB to store Objects
- [SQL](https://github.com/wazofski/storz/tree/main/sql) store - uses a SQL database connection for storage

### Functional Layer
Functional modules require existing Stores to operate.
These modules are meant to enhance the functionality of an existing store by composing itself with another Store.
Caching layer can be added to a Store and then wrapped into another layer of React that adds validation logic object changes.

- [Cache](https://github.com/wazofski/storz/tree/main/cache) store - simple cachingg mechanism using an existing Store
- [Route](https://github.com/wazofski/storz/tree/main/route) store - mapping between types and Stores is used to route requets
- [React](https://github.com/wazofski/storz/tree/main/react) store - react to object changes before they get submitted

### REST
- [Server](https://github.com/wazofski/storz/tree/main/rest)
- [Client](https://github.com/wazofski/storz/tree/main/client) store

### Utility
- [Browser](https://github.com/wazofski/storz/tree/main/browser)
