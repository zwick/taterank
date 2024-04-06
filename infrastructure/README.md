# Taterank Infrastructure

Taterank uses Terraform to deploy infrastructure to AWS. In order to deploy this infrastructure,
you will require both [Terraform](https://developer.hashicorp.com/terraform/install) and the 
[AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) with a profile 
setup with sufficient permissions to create resources in the AWS account.

## Usage

### Create and Modify Infrastructure
To instantiate Terraform and download all required dependencies, use

```shell
terraform init
```

Next, run a plan against the current infrastructure. 

```shell
terraform plan
```

If there are any changes to be made based on the current infrastucture, terraform will display what changes it intends
to make. Proceed to the next step to apply these changes to the infrastructure.

```shell
terraform apply
```

This will apply any required changes to the deployed AWS resources.

### Destroy Infrastructure

This infrastructure can be easily tore down. This is useful for temporary environments for testing and save costs. To 
destroy your infrastructure, run

```shell
terraform destroy
```

Be very careful. Running this command will result in the loss of all data in the DynamoDB instance.
