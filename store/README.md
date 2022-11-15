# Store Definitions

The Store common interface defined below exposes
five main operations over Objects allowing various
options like pagination, filtering and sorting.

**Store common interface**

```
  Get(context.Context, ObjectIdentity, ...GetOption) (Object, error)
  List(context.Context, ObjectIdentity, ...ListOption) (ObjectList, error)
  Create(context.Context, Object, ...CreateOption) (Object, error)
  Delete(context.Context, ObjectIdentity, ...DeleteOption) error
  Update(context.Context, ObjectIdentity, Object, ...UpdateOption) (Object, error)
```

