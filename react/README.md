# React

React store exposes a way to attach callbacks to object 
actions associated with an underlying store

## Usage

```
store := store.New(
    generated.Schema(),
    react.ReactFactory(underlying_store))
```
