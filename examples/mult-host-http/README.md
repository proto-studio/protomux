# Multi-Host HTTP Example

This is an example of using Proto Mux to route traffic for multiple routes on a different domain.

This demonstrates:

- Setting up the routes.
- Using route parameters.
- Using domain parameters.

The app has three endpoints:

- `GET /doc/{id}` gets an existing document or returns `404`
- `PUT /doc/{id}` stores the request body under the document ID.
- `DELETE /doc/{id}` deletes the document or returns `404`.

You can also try other urls or different HTTP methods to see how you get different responses.

Proto Mux will automatically respond with a `404` for missing routes and a `405` if there is a registered route but that method isn't registered.

This example works exactly like the "single host http" example except that you can have more than one "database" on the same server.

This app listens to `*.example.com`. Any other domain will return `404`. The subdomain portion of the URL is the "database name" for the request.

You can test this with something like `curl` or you can setup hostnames in your hosts file for each subdomain.

## Using Curl

For example, using curl, you can create a document with the key `hello` in the `db1` subdomain using:

```bash
curl -X PUT -d "world" -H "Host: db1.example.com" localhost:8080/docs/hello
```

And retrieve it using:

```bash
curl  -H "Host: db1.example.com" localhost:8080/docs/hello
```

Or you can try it with a different host, which will return `404` since this record only exists in the `db1` subdomain:

```bash
curl  -H "Host: db2.example.com" localhost:8080/docs/hello
```