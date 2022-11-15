# Router

Router store allows mapping object types to different Stores

## Usage

```
store := store.New(
    generated.Schema(),
    router.Factory(deault_store,
        router.Mapping("type1", store1),
        router.Mapping("type2", store2)))
```
