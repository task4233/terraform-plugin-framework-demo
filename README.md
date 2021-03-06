## Terraform-Plugin-Framework-Demo
This is an unofficial demo for [Terraform-Plugin-Framework](https://www.terraform.io/docs/plugin/framework/) which is a new SDK under active development.
You might be able to get further information in [Official GitHub Repository](https://github.com/hashicorp/terraform-plugin-framework).

※ This plugin needs Terraform CLI v1.0 or later.
ref: [Terraform Plugin Framework](https://www.terraform.io/plugin/framework)

## How to Use
1. install plugin locally

```bash
$ make install
```

2. terraform init

```bash
cd examples
rm -rf terraform.*
terraform init
```

3. terraform apply(Create)

```bash
echo "running demo server for serving data"
go run ../demo-server/main.go
echo "invoke `terraform apply` on another shell"
terraform apply
```

4. terraform apply(Read)

```bash
terraform apply
```

5. terraform apply(Update)

```bash
echo "edit configuration in main.tf"
vim main.tf
terraform apply
```

6. terraform destroy(Delete)

```bash
terraform destroy
```

## Useful Information
- [Plugin Development](https://www.terraform.io/docs/extend/index.html)
- [Implement Create and Read with the Terraform Plugin Framework](https://learn.hashicorp.com/tutorials/terraform/plugin-framework-create)
- [terraform-provider-hashicups-pf](https://github.com/hashicorp/terraform-provider-hashicups-pf)
- [hashicups-client-go](https://github.com/hashicorp-demoapp/hashicups-client-go)
