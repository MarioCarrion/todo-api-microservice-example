# OpenAPI 3 / Swagger

* [OpenAPI 3 Specification (3.0.0)](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.0.md)

## Tools

* [github.com/getkin/kin-openapi/openapi3](https://pkg.go.dev/github.com/getkin/kin-openapi/openapi3#section-documentation): Package openapi3 parses and writes OpenAPI 3 specifications.
* [Swagger Editor](https://editor.swagger.io/)
* [Swagger UI](https://github.com/swagger-api/swagger-ui), local copy is in [`cmd/rest-server/static/swagger-ui`](../cmd/rest-server/static/swagger-ui).
    * Local demo: [http://0.0.0.0:9234/static/swagger-ui/](http://0.0.0.0:9234/static/swagger-ui/)

### OpenAPI Codegen 3.X

For Go the types representing the REST API in [`.`](../.): [`oapi-codegen`](https://github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen) is used for generating them.

For other languages you may want to use their corresponding equivalent, for example for Ruby you can use `swaggerapi/swagger-codegen-cli-v3`:

```
docker run --rm -v ${PWD}:/gen swaggerapi/swagger-codegen-cli-v3:3.0.25 generate \
  --verbose \
  --input-spec http://host.docker.internal:9234/openapi3.yaml \
  --lang ruby \
  --output /gen
```
