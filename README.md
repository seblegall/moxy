# Moxy

A proxy AND mock server.

Moxy let you mock api endpoint you need to be mock and will proxy all other endpoint.


## Example

You are working on a frontend application which consume 2 endpoints :

* `GET /myressources`
* `POST /myressources`

In order to test all case, you need the GET endpoint to return a specific value (null for example).

You can tell Moxy to mock the `GET /myressources` endpoint in order to get what you need or expect. All other endpoint will be called exactly as if 
you were calling them.


## Usage

```
moxy -backend=http://mydomain.com -port=8080
```

The `backend` flag is mandatory. The `port` flag has a default value (8080).
