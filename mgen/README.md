# Model GENerator

The model is a collection of yaml files holding Store Object and Structure definitions. 

## Model
### Objects
Modeled Object describes the high-level abstraction used to manipulate all data via the Store inteface.

The Object is structured as follows
- Metadata
    - Kind information
    - Primary key
    - Framework assigned identitier
    - Object manipulation timestamps (create, update...)
- Spec (External) - any Structure, to be managed through external REST APIs
- Status (Internal) - to be managed by internal service code (React callbacks)

The Object contains 

### Structures
Structures are named collections of typed properties.

Supported property types
- Golang standard types
    - string
    - int
    - float
    - bool
- Other Structures (nesting)
- String-keyed maps
    - map[string]int (string, float...)
    - map[string]Struct
- Slices
    - []int (string, float...)
    - []Struct


## Generated

Import the "generated" package to access Object interfaces and Schema.

Use <object>Factory() functions to create Model specific mutable Objects
which contain Object
- Metadata
- Spec / Status
    - Property Getters/Setters
