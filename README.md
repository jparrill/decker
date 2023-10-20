# Decker
Decker stands for "Disconnected Checker" and it's a tool which helps you out validating Openshift disconnected deployments in multiple ways. Among other features you can:

- Verify an OCP Pull Secret format
- Verify a container registry access and/or a container registry image
- Generate ICSP/IDMS manifests files to be applied in Management clusters
- Validate a OCP Release payload
- Diagnose a OCP disconnected deployment


## Verify

This commands basically verifies a Pull Secret input file:

```bash
decker verify pull-secret --authfile pkg/verify/test.json
```

```
Verifying pullsecret: pkg/verify/test.json
✔︎ - Read input file
✔︎ - Unmarshal JSON file
```

If you wanna inspect the registries inside of that pull secret:

```bash
decker verify pull-secret --authfile pkg/verify/test.json --inspect
```

```bash
Verifying pullsecret: pkg/verify/test.json
✔︎ - Read input file
✔︎ - Unmarshal JSON file

RegistryName: registry.lab04:5000
⨯ - Registry Credentials
    • Error: No authentication provided
⨯ - Registry Authentication
    • Error: Error login into destination registry

RegistryName: registry.lab01:5000
✔︎ - Registry Credentials
⨯ - Registry Authentication
    • Error: Error login into destination registry

RegistryName: registry.lab02:5000
✔︎ - Registry Credentials
⨯ - Registry Authentication
    • Error: Error login into destination registry

RegistryName: registry.lab03:5000
✔︎ - Registry Credentials
⨯ - Registry Authentication
    • Error: Error login into destination registry
```