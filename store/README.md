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

## Usage

```
// Given a Store
var store_ store.Store = ...

// Initialize the object
world := generated.WorldFactory()
world.Spec().SetName("abc")

// Create the object on
world, err = store_.Create(ctx, world)

world.Spec().SetDescription("abc")

// Update the object
world, err = clt.Update(ctx, world.Metadata().Identity(), world)

// Delete the object
err = clt.Delete(ctx, world.Metadata().Identity())

// Get the object by Identity or PKey
world, err = store_.Get(ctx, generated.WorldIdentity("abc"))

// List the World objects
world_list, err = store_.List(ctx, generated.WorldIdentity(""))
```