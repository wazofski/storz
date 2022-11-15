# STORZ

Storz is an Object Store framework built in golang.

It features a simple object modeling language used to [generate ](https://github.com/wazofski/storz/tree/main/mgen)
the golang object class meta.

The generated code contains object and structure classes used to interact 
with the Store interface that most storz modules expose/implement.


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
Most modules are implementations of the [Store](https://github.com/wazofski/storz/tree/main/store) interface.


### Persistification
- Memory store
- SQL store
  - mySQL
  - sqlite
- Mongo store


### Functional
- Cached store
- [React](https://github.com/wazofski/storz/tree/main/react) store

  
### REST
- [Server](https://github.com/wazofski/storz/tree/main/rest)
- [Client](https://github.com/wazofski/storz/tree/main/client) (store)


### Other
- Object Logger (store)
- Object Browser

