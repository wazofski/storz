<p align="center">
<img src="logo.png" width="300" alt="storz" />
</p>

<!-- ![storz](https://github.com/wazofski/storz/blob/main/logo.png?raw=true) -->

**storz** is an *object store framework* built in golang. It consists of a set of modules implementing the [Store](https://github.com/wazofski/storz/tree/main/store) interface and features a simple [object modeling language](https://github.com/wazofski/storz/tree/main/mgen) used to generate golang object class meta for interacting with `Store` APIs.

**storz** modules provide functionality to store, modify and retrieve modeled objects from various sources. Such modules can be composed together to chain Store functionality into more complex logical modules. Combining modules allows handling object changes and manipulating data in complex ways *within or across* services making multi-level server complexity achievable with ease. 

Example:
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

## Features
**storz** modules all implement or expose the [Store](https://github.com/wazofski/storz/tree/main/store) interface and are fully compatible with each other, so comoposing modules is allowed in any combination.

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
