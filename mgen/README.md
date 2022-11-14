# Model GENerator

The model is a collection of yaml files holding Store Object and Structure definitions. 

## Objects
Modeled Object describes the high-level abstraction used to manipulate all data via the Store inteface.

The Object is structured as follows
- Metadata
    - Kind information
    - Primary key
    - Framework assigned identitier
    - Object manipulation timestamps (create, update...)
- External
    - Can be assigned any Structure, can be modified by the external APIs
- Internal
    - Can be assigned any Structure, cannot be modified by external APIs


## Structures
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

