# OpenAPI 3

* [OpenAPI 3 Specification (3.0.0)](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.0.md)

## Tools

* [oapi-codegen/oapi-codegen](https://github.com/oapi-codegen/oapi-codegen): OpenAPI 3 code generator.
* [Swagger Editor](https://editor.swagger.io/)
* [Swagger UI](https://github.com/swagger-api/swagger-ui), local copy is in [`cmd/rest-server/static/swagger-ui`](../cmd/rest-server/static/swagger-ui).
    * Local demo: [http://localhost:9234/static/swagger-ui/](http://localhost:9234/static/swagger-ui/)

### OpenAPI 3 type-safe generated code

This project uses [`oapi-codegen`](https://github.com/oapi-codegen/oapi-codegen/tree/main/cmd/oapi-codegen) to generate both the server and client models using a [Top-Bottom](https://youtu.be/ErA92edMta8) approach where the [OpenAPI specification](../openapi/openapi3.yaml) is written first.

* [Server configuration](../internal/rest/oapi-codegen.server.yaml)
* [Client configuration](../oapi-codegen.client.yaml)
