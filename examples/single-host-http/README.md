# Single Host HTTP Example

This is an example of using Proto Mux to route traffic for multiple routes on a single domain.

This demonstrates:

- Setting up the routes.
- Registering route parameters.

The app has three endpoints:

- `GET /doc/{id}` gets an existing document or returns `404`
- `PUT /doc/{id}` stores the request body under the document ID.
- `DELETE /doc/{id}` deletes the document or returns `404`.

You can also try other urls or different HTTP methods to see how you get different responses.

Proto Mux will automatically respond with a `404` for missing routes and a `405` if there is a registered route but that method isn't registered.

## Using Curl

For example, using curl, you can create a document with the key `hello` using:

```bash
curl -X PUT -d "world" localhost:8080/docs/hello
```

And retrieve it using:

```bash
curl localhost:8080/docs/hello
```
