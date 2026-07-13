AWS Serverless Web App

What this template creates:
- 1 S3 bucket for the static frontend
- 1 CloudFront distribution for secure global frontend delivery
- 1 Lambda function for the backend API
- 1 HTTP API Gateway in front of the Lambda function
- 1 DynamoDB table for application data
- 1 Lambda execution role and DynamoDB access policy

Recommended use:
- Good for startups that want the lowest practical cloud cost
- Works well for MVPs, internal tools, dashboards, and simple customer-facing apps

Cost notes:
- No VPC, NAT gateway, or RDS by default
- DynamoDB on-demand keeps the database layer simple and low-cost
- CloudFront reduces frontend delivery cost and improves performance

Customization ideas:
- Replace the sample frontend page with your own app build output
- Add Cognito or another auth layer if you need sign-in
- Swap the sample Lambda handler with your business logic
