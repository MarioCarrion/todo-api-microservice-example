# Secure Configuration

## HashiCorp Vault

Used as repository for retrieving secrets for configuration values.

Please review the `vault` **services** in [compose.yml](../compose.yml), the code to read secure values
is in the [vault](../internal/envvar/vault) package.
