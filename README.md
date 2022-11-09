# storz
Storz is an Object Store framework built in golang.
It features a simple object modeling language used to generate the golang object class meta.
The generated code contains object and structure classes used to interact with the Store 
interface that most storz modules expose/implement.

## Store common interface

Most modules is a Store implementation.

```
  Get(context.Context, ObjectIdentity, ...GetOption) (Object, error)
  List(context.Context, ObjectIdentity, ...ListOption) (ObjectList, error)
  Create(context.Context, Object, ...CreateOption) (Object, error)
  Delete(context.Context, ObjectIdentity, ...DeleteOption) error
  Update(context.Context, ObjectIdentity, Object, ...UpdateOption) (Object, error)
```


# Getting started

1. Install the storz module

```
  go get github.com/wazofski/storz.git
  go install github.com/wazofski/storz.git
```

2. Initialize your storz project

```
  storz init <project>
```

This will create the <project> directory containing your go module, a sample model and the main.go source file.

3. Generate the class meta

In your <project> directory, run

```
  go generate
```





## rest

### server

### client store



## react store





## memory store

## cached store

## sql store

## mongodb store

# other




## visore


## logging





