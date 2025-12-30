# APSchool - AWS Deployment Guide

This document describes the AWS infrastructure configuration for the APSchool production environment.

## Architecture Overview

```
                    CloudFront (HTTPS)
                    d3ow6eel2jvdfj.cloudfront.net
                          │
            ┌─────────────┴─────────────┐
            │                           │
       /* (default)                  /api/*
            │                           │
            ▼                           ▼
     S3 Bucket                    EC2 Instance
   (Frontend Angular)            (Backend Go:8080)
                                        │
                                        ▼
                                 RDS PostgreSQL
                                   (port 5432)
```

## AWS Services

### 1. Security Groups

#### `apschool-ec2-sg`
Security group for the EC2 backend instance.

| Type | Port | Source | Description |
|------|------|--------|-------------|
| SSH | 22 | 0.0.0.0/0 | SSH access |
| HTTP | 80 | 0.0.0.0/0 | HTTP traffic from CloudFront |

#### `apschool-rds-sg`
Security group for the RDS PostgreSQL database.

| Type | Port | Source | Description |
|------|------|--------|-------------|
| PostgreSQL | 5432 | apschool-ec2-sg | Database access from EC2 only |

---

### 2. RDS PostgreSQL

| Property | Value |
|----------|-------|
| Instance Identifier | `apschool-db` |
| Engine | PostgreSQL 18.1 |
| Instance Class | `db.t4g.micro` (Free Tier) |
| Storage | 20 GiB gp2 |
| Endpoint | `apschool-db.cgpmss6621un.us-east-1.rds.amazonaws.com` |
| Port | 5432 |
| Database Name | `apschool` |
| Master Username | `apschool_admin` |
| Public Access | No |
| Security Group | `apschool-rds-sg` |
| SSL Mode | `require` |

---

### 3. EC2 Instance

| Property | Value |
|----------|-------|
| Name | `apschool-backend` |
| Instance ID | `i-0c9af5529f1d9ed97` |
| Instance Type | `t3.micro` (Free Tier) |
| AMI | Ubuntu Server 24.04 LTS |
| Public IP | `3.238.148.89` |
| Public DNS | `ec2-3-238-148-89.compute-1.amazonaws.com` |
| Key Pair | `apschool-key` |
| Security Group | `apschool-ec2-sg` |
| Storage | 8 GiB gp3 |

#### Installed Software
- Docker 29.1.3

#### SSH Access
```bash
ssh -i ~/.ssh/apschool-key.pem ubuntu@3.238.148.89
```

---

### 4. S3 Bucket

| Property | Value |
|----------|-------|
| Bucket Name | `apschool-frontend-213590135603` |
| Region | us-east-1 |
| Bucket Type | General Purpose |
| Public Access | Blocked (accessed via CloudFront OAC) |
| Versioning | Disabled |
| Encryption | SSE-S3 |

---

### 5. CloudFront Distribution

| Property | Value |
|----------|-------|
| Distribution ID | `E372A787F4U1GQ` |
| Domain Name | `d3ow6eel2jvdfj.cloudfront.net` |
| Price Class | Use only North America and Europe |
| Default Root Object | `index.html` |
| Viewer Protocol Policy | Redirect HTTP to HTTPS |

#### Origins

| Name | Type | Domain |
|------|------|--------|
| S3 Origin | S3 Bucket | `apschool-frontend-213590135603.s3.us-east-1.amazonaws.com` |
| EC2 Origin | Custom | `ec2-3-238-148-89.compute-1.amazonaws.com` |

#### Behaviors

| Path Pattern | Origin | Methods | Cache Policy |
|--------------|--------|---------|--------------|
| `/api/*` | EC2 Origin | GET, HEAD, OPTIONS, PUT, POST, PATCH, DELETE | CachingDisabled |
| `/*` (default) | S3 Origin | GET, HEAD | CachingOptimized |

---

## GitHub OAuth App (Production)

| Property | Value |
|----------|-------|
| Application Name | APSchool Production |
| Client ID | `Ov23likDtafYtW94THVx` |
| Homepage URL | `https://d3ow6eel2jvdfj.cloudfront.net` |
| Callback URL | `https://d3ow6eel2jvdfj.cloudfront.net/api/auth/github/callback` |

---

## Environment Variables (Production)

The `.env.prod` file should be created on the EC2 instance at `/home/ubuntu/.env.prod`:

```env
# Database (RDS PostgreSQL)
APSCHOOL_DB_HOST=apschool-db.cgpmss6621un.us-east-1.rds.amazonaws.com
APSCHOOL_DB_PORT=5432
APSCHOOL_DB_DATABASE=apschool
APSCHOOL_DB_USERNAME=apschool_admin
APSCHOOL_DB_PASSWORD=<YOUR_RDS_PASSWORD>
APSCHOOL_DB_SCHEMA=public
APSCHOOL_DB_SSLMODE=require

# Server
PORT=8080

# JWT
JWT_SECRET=<GENERATE_WITH_openssl_rand_-base64_32>

# GitHub OAuth (Production)
GITHUB_CLIENT_ID=Ov23likDtafYtW94THVx
GITHUB_CLIENT_SECRET=<YOUR_CLIENT_SECRET>
GITHUB_REDIRECT_URI=https://d3ow6eel2jvdfj.cloudfront.net/api/auth/github/callback

# Frontend
FRONTEND_URL=https://d3ow6eel2jvdfj.cloudfront.net
```

---

## Deployment Commands

### Backend (EC2)

```bash
# SSH into EC2
ssh -i ~/.ssh/apschool-key.pem ubuntu@3.238.148.89

# Navigate to project
cd ~/apschool

# Pull latest changes
git pull origin main

# Build Docker image
cd server
docker build -t apschool-backend .

# Stop existing container (if running)
docker stop apschool-backend 2>/dev/null
docker rm apschool-backend 2>/dev/null

# Run container
docker run -d --name apschool-backend \
  -p 80:8080 \
  --env-file /home/ubuntu/.env.prod \
  --restart unless-stopped \
  apschool-backend

# Run migrations (first time or when migrations change)
docker exec apschool-backend goose -dir /app/migrations postgres \
  "postgres://apschool_admin:<PASSWORD>@apschool-db.cgpmss6621un.us-east-1.rds.amazonaws.com:5432/apschool?sslmode=require" up

# Run seed (first time or when challenges change)
docker exec apschool-backend ./seed

# Check logs
docker logs -f apschool-backend
```

### Frontend (S3)

```bash
# From local machine, in web/ directory
cd web

# Build for production
bun run build

# Sync to S3
aws s3 sync dist/web/browser s3://apschool-frontend-213590135603 --delete

# Invalidate CloudFront cache (optional)
aws cloudfront create-invalidation \
  --distribution-id E372A787F4U1GQ \
  --paths "/*"
```

---

## URLs

| Environment | URL |
|-------------|-----|
| Production | https://d3ow6eel2jvdfj.cloudfront.net |
| API | https://d3ow6eel2jvdfj.cloudfront.net/api |

---

## Free Tier Limits

| Service | Free Tier Allowance | Our Usage |
|---------|---------------------|-----------|
| EC2 (t3.micro) | 750 hours/month | 1 instance |
| RDS (db.t4g.micro) | 750 hours/month | 1 instance |
| S3 | 5 GB storage | ~10 MB |
| CloudFront | 1 TB transfer, 10M requests | Minimal |
| Data Transfer | 100 GB/month | Minimal |

---

## Important Notes

1. **EC2 Public IP**: The public IP may change if the instance is stopped/started. Consider using an Elastic IP for a static address.

2. **RDS Access**: The database is not publicly accessible. It can only be accessed from the EC2 instance.

3. **SSL/TLS**: All public traffic is encrypted via CloudFront. The connection between CloudFront and EC2 is HTTP (internal AWS network).

4. **Secrets**: Never commit `.env.prod` or any credentials to the repository.
