# AWS Service Group config

Dynamically modifying rules in Service Groups using the AWS Go SDK.

This code can be run in a pipeline or in a Lambda as an automation.

Create the base infrastructure:

```sh
terraform -chdir="infra" init
terraform -chdir="infra" apply -auto-approve
```

The file `cirds.yaml` contains the SG rule configuration.

Get and run:

```sh
go get
go run .
```

Check the SG for the modifications.

After you're done, remove the resources:

```sh
terraform -chdir="infra" destroy -auto-approve
```
