# Secure Configuration

## HashiCorp Vault

Used as repository for retrieving secrets for configuration values.

```
docker run \
  -d \
  --cap-add=IPC_LOCK \
  -e 'VAULT_DEV_ROOT_TOKEN_ID=myroot' \
  -e 'VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:8300' \
  -p 8300:8300 \
  vault:1.6.2
```
