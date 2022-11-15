# REST

REST package provides Store exposing REST API server

## Usage

```
srv := rest.Server(generated.Schema(), store_to_expose)

// use cancel function to stop server
cancel = srv.Listen(port) // does not block
```
