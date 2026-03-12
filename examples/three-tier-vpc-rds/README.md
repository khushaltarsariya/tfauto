# Three-Tier VPC + RDS Example

This example shows a simple two-step project flow:

1. Create the network project from `aws-three-tier-vpc`
2. Create the database project from `aws-rds-postgres`

## 1. Create the VPC project

```bash
tfauto init --template aws-three-tier-vpc --target ./network
```

Suggested values:
- keep the default private DB subnet ranges
- use `single_nat_gateway = true` if you want a lower-cost dev setup

## 2. Create the database project

```bash
tfauto init --template aws-rds-postgres --target ./database
```

Wire these values from the VPC outputs into the database project:
- `vpc_id`
- `db_subnet_ids`
- `app_security_group_id` into `allowed_security_group_ids`

## Suggested flow

```bash
tfauto validate --path ./network
tfauto plan --path ./network --out network.tfplan
tfauto apply --path ./network --plan network.tfplan

tfauto validate --path ./database
tfauto plan --path ./database --out database.tfplan
tfauto apply --path ./database --plan database.tfplan
```
