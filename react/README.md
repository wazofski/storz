# React

React store exposes a way to attach callbacks to object 
actions associated with an underlying store

## Usage

```
func WorldCreateCb(store.Object, store.Store) error {
    // ...
    return nil
}

func WorldDeleteCb(store.Object, store.Store) error {
    // ...
    return nil
}

store := store.New(
    generated.Schema(),
    react.ReactFactory(underlying_store,
        react.Subscribe("World", react.ActionCreate, WorldCreateCb),
        react.Subscribe("World", react.ActionDelete, WorldDeleteCb),
    ))
```
