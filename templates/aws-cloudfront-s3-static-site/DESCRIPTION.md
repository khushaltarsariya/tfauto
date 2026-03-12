AWS CloudFront + S3 Static Site

What this template creates:
- 1 private S3 bucket for site assets
- 1 Origin Access Control for CloudFront
- 1 CloudFront distribution
- 1 bucket policy allowing only CloudFront to read objects
- Static index and error document settings

Recommended use:
- Use this instead of direct public S3 website hosting for production-style static sites
- Pair with Route 53 and ACM if you want a custom domain

Notes:
- This template keeps the S3 bucket private
- CloudFront is the public entry point for the site
