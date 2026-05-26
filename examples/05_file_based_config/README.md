# File-Based Configuration Example

This example demonstrates using the `file` tag to load configuration from files.

## Feature: File Reading (`file` tag)

The `file` tag tells envy to treat the environment variable value as a file path and read the file content.

```go
TLSCert string `env:"TLS_CERT,file"`
```

## Why Use File Reading?

### 1. **Security**
Keep sensitive data (certificates, keys, credentials) in separate files with restricted permissions instead of environment variables.

```bash
# Instead of: export DB_PASSWORD=supersecret
# Use environment variable pointing to file:
export DB_PASSWORD_FILE=/run/secrets/db_password
# Then read: PrivateKey string `env:"DB_PASSWORD_FILE,file"`
```

### 2. **Large Configuration**
Multi-line configuration files are easier to manage than escaped environment variables.

```bash
export DB_CONFIG=/etc/config/database.conf
# File contains:
# [database]
# host = localhost
# port = 5432
# ...
```

### 3. **Docker/Kubernetes Integration**
Docker Secrets and Kubernetes Secrets often provide files in a mounted directory.

### 4. **Dynamic Configuration**
Update configuration files without restarting the application (if you implement reload logic).

---

## Usage Examples

### Basic File Reading

```go
type Config struct {
    CertFile string `env:"TLS_CERT,file"`
}
```

**Environment setup:**
```bash
export TLS_CERT=/etc/ssl/certs/server.crt
# envy reads the file content into CertFile field
```

### With Defaults

```go
type Config struct {
    // If env var not set, uses the literal default string (not as file path)
    PrivateKey string `env:"PRIVATE_KEY,default=-----BEGIN PRIVATE KEY-----\nMIIEvQ..."`
}
```

### Multiple File-Based Fields

```go
type Config struct {
    TLSCert    string `env:"TLS_CERT_FILE,file,required"`
    PrivateKey string `env:"TLS_KEY_FILE,file,required"`
    ConfigData string `env:"APP_CONFIG_FILE,file,default=/etc/app/default.conf"`
}
```

---

## Docker Example

```dockerfile
FROM golang:1.21

WORKDIR /app
COPY . .

# Docker will inject secrets as files
# You can reference them:
ENV TLS_CERT=/run/secrets/tls_cert
ENV DB_PASSWORD_FILE=/run/secrets/db_password

CMD ["go", "run", "main.go"]
```

```bash
# Run with Docker secrets
docker run \
  --secret tls_cert \
  --secret db_password \
  -e TLS_CERT=/run/secrets/tls_cert \
  -e DB_PASSWORD_FILE=/run/secrets/db_password \
  myapp
```

---

## Kubernetes Example

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: myapp
spec:
  containers:
  - name: app
    image: myapp:latest
    env:
    - name: TLS_CERT_FILE
      value: /var/run/secrets/tls/tls.crt
    - name: TLS_KEY_FILE
      value: /var/run/secrets/tls/tls.key
    volumeMounts:
    - name: tls-secret
      mountPath: /var/run/secrets/tls
      readOnly: true
  volumes:
  - name: tls-secret
    secret:
      secretName: app-tls
```

---

## Running the Example

```bash
cd examples/05_file_based_config

# Create test files
mkdir -p /tmp/config
echo "-----BEGIN CERTIFICATE-----" > /tmp/config/cert.pem
echo "MIIDXTCCAkWgAwIBAgI..." >> /tmp/config/cert.pem
echo "-----END CERTIFICATE-----" >> /tmp/config/cert.pem

echo "-----BEGIN PRIVATE KEY-----" > /tmp/config/key.pem
echo "MIIEvQIBADANBgkqhkiG9w0..." >> /tmp/config/key.pem
echo "-----END PRIVATE KEY-----" >> /tmp/config/key.pem

# Set environment variables
export APP_NAME=myapp
export TLS_CERT=/tmp/config/cert.pem
export PRIVATE_KEY=/tmp/config/key.pem

go run main.go
```

---

## Error Handling

If the file doesn't exist, envy returns a validation error:

```
1 env error(s):
  - TLSCert (TLS_CERT): failed to read file "/nonexistent/cert.pem": open /nonexistent/cert.pem: no such file or directory
```

---

## Best Practices

1. **Use `required` for critical files:**
   ```go
   TLSCert string `env:"TLS_CERT,file,required"`
   ```

2. **Provide meaningful defaults:**
   ```go
   LogConfig string `env:"LOG_CONFIG_FILE,file,default=/etc/app/logging.conf"`
   ```

3. **Use restricted file permissions:**
   ```bash
   chmod 600 /var/secrets/db_password
   ```

4. **Document file paths in comments:**
   ```go
   // TLSCert: Path to TLS certificate file (typically /etc/ssl/certs/server.crt)
   TLSCert string `env:"TLS_CERT,file,required"`
   ```

5. **Consider file size for large files** - small files in env vars, large files in files

---

## Key Takeaways

- Use `file` tag for security-sensitive data
- Works well with Docker Secrets and Kubernetes Secrets
- Supports defaults with literal strings (not file paths)
- Errors are collected like other validation errors
- Combine with `required` and `notEmpty` for validation
