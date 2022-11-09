# STORZ

Storz is an Object Store framework built in golang.
It features a simple object modeling language used to generate the golang object class meta.
The generated code contains object and structure classes used to interact with the Store 
interface that most storz modules expose/implement.


# Getting started

1. Install the storz module

```
  go get github.com/wazofski/storz.git
  go install github.com/wazofski/storz.git
```

2. Initialize your storz project

```
  storz init [project]
```

This will create the [project] directory containing your go module, a sample model and the main.go source file.

3. Generate the class meta

In your [project] directory, run

```
  go generate
```

Don't forget to rerun go generate when making changes to your model YAML files.
  
  
4. Build and run your code
```
  go build
  ./[project]
```

# Modules

## Store common interface

Most modules are implementations of the following common interface.
```
  Get(context.Context, ObjectIdentity, ...GetOption) (Object, error)
  List(context.Context, ObjectIdentity, ...ListOption) (ObjectList, error)
  Create(context.Context, Object, ...CreateOption) (Object, error)
  Delete(context.Context, ObjectIdentity, ...DeleteOption) error
  Update(context.Context, ObjectIdentity, Object, ...UpdateOption) (Object, error)
```


## Persistification modules
### In-memory store

### SQL store

### Mongo DB store
  
  
## Functional modules
### Cache store
### React store
React provides a way to attach callbacks to object actions associated with an underlying store

  
## REST modules
### REST Server
### REST Client store


  
## Other useful modules
### Detailed Logger store
### Object Browser server

  
