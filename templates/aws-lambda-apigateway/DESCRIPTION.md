AWS Lambda + API Gateway

What this template creates:
- 1 Lambda execution role
- 1 Lambda function from inline zip content
- 1 HTTP API Gateway
- 1 default stage with auto-deploy
- Integration and route wiring between API Gateway and Lambda
- Permission allowing API Gateway to invoke the function

Recommended use:
- Use this as a starter for lightweight APIs, webhooks, and internal platform endpoints
- Replace the inline handler code with your own zip artifact for real workloads

Notes:
- The default function returns a simple JSON response
- Switch to artifact-based deployment for larger functions
