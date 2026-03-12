# Three-Tier VPC + ALB/ASG Web App Example

This example combines:
- `aws-three-tier-vpc`
- `aws-alb-asg-webapp`

## 1. Create the network project

```bash
tfauto init --template aws-three-tier-vpc --target ./network
```

## 2. Create the web application project

```bash
tfauto init --template aws-alb-asg-webapp --target ./webapp
```

Wire these values from the network outputs into the web app project:
- `vpc_id`
- `public_subnet_ids` into `alb_subnet_ids`
- `app_subnet_ids` into `instance_subnet_ids`

## Recommended plan/apply flow

```bash
tfauto plan --path ./network --out network.tfplan
tfauto apply --path ./network --plan network.tfplan

tfauto plan --path ./webapp --out webapp.tfplan
tfauto apply --path ./webapp --plan webapp.tfplan
```
