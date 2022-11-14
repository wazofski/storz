# Store Definitions

**Store common interface**

```
  Get(context.Context, ObjectIdentity, ...GetOption) (Object, error)
  List(context.Context, ObjectIdentity, ...ListOption) (ObjectList, error)
  Create(context.Context, Object, ...CreateOption) (Object, error)
  Delete(context.Context, ObjectIdentity, ...DeleteOption) error
  Update(context.Context, ObjectIdentity, Object, ...UpdateOption) (Object, error)
```

