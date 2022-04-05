# All Triumph Challenge
`base64.go` is a tiny program that generates `favicon.go` source code, so no build environment is required.

The favicon and small ATC logo are used by the HTML result output. These 2 files are then included into the HTML results source, so no external files need to be distributed along with the HTML file.

### How to generate `favicon.go`
After modifying images `/cmd/favicon.ico` or `/cmd/logo.png` then `/favicon.go` needs to be regenerated to see future HTML results updated. This can be done by executing the following commands:
```shell
cd TriumphChallenge/cmd/
go run base64.go
gofmt -w .
mv favicon.go ../favicon.go
```