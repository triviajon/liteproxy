## liteproxy

liteproxy reduces the size of HTTP responses by intercepting requests, rewriting content, and responding to users.

### Dependencies

- GO_VERSION=1.26.1
- TERRAFORM_VERSION=1.14.7

### Building

Build the processor image:

```
source versions.env && docker build --build-arg GO_VERSION=$GO_VERSION -t liteproxy-processor:dev .
```
