# Lambda + API Gateway Example

This example uses the `aws-lambda-apigateway` template for a minimal HTTP API.

## Create the project

```bash
tfauto init --template aws-lambda-apigateway --target ./api
```

## Validate and deploy

```bash
tfauto validate --path ./api
tfauto plan --path ./api --out api.tfplan
tfauto apply --path ./api --plan api.tfplan
```

## Customize

- replace the inline handler logic with your own implementation
- add environment variables through `terraform.tfvars`
- front the API with a custom domain if needed
