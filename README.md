# defect-detect

```
# download curl and install trivy
apk add curl
curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin v0.52.2

# download test sbom and scan
curl "https://github.com/CycloneDX/bom-examples/blob/master/SBOM/keycloak-10.0.2/bom.json?raw=true" -L > bom.json
trivy sbom bom.json -f json -o test.json
```