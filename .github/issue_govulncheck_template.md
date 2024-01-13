---
title: "Vulnerabilities found"
labels: security, govulncheck
---
[Govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) has found some issues in the latest GitHub Action run: <{{ env.ACTION_FULL_URL }}>.

Full output:

````console
$ govulncheck ./...
{{ env.GOVULNCHECK_OUTPUT }}
````
