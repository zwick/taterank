# Taterank

The best spud ranking app to ever exist

## Development

To run the app locally, install the following:

- [Go](https://go.dev/)
- [Docker](https://www.docker.com/products/docker-desktop/)

Execute the `bootstrap` script under `/bin`.

Run `docker compose up`.

## Deployment

### API
To deploy the API, you need to have the following:

- Go
- AWS CLI
- Serverless framework
- Terraform

In the `infrastructure` directory, run the following command to deploy the infrastructure:

```bash
terraform apply
```
This will create all necessary AWS resources including the DynamoDB table, IAM permissions and role required for the 
Lambda functions.

Next, from the `server` directory build the Go app into a file called `bootstrap`.

```bash
GOARCH=amd64 GOOS=linux go build -o bootstrap ./cmd/api
```

Then once you have the `bootstrap` file, you can deploy the API using the following command:

```bash
serverless deploy
```

This will deploy the app to AWS and return the API Gateway URL in the output.
