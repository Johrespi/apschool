# APSchool

<p align="center">
  <img src="https://img.shields.io/badge/Angular-21.0.6-dd0031" alt="Angular">
  <img src="https://img.shields.io/badge/Go-1.25.5-00add8" alt="Go">
  <img src="https://img.shields.io/badge/Terraform-1.16+-844FBA" alt="Terraform">
  <img src="https://img.shields.io/badge/AWS-Cloud-232F3E" alt="AWS">
</p>

[English](README.md) | [Español](README.es.md)

An interactive web platform designed for computer science and engineering students (ESPOL) to master Python programming through hands-on algorithmic challenges, with code execution running entirely client-side via Pyodide (WebAssembly).

---

## The Problem

Learning programming in foundational university courses presents significant hurdles: complex local environment setups for beginners, lack of immediate automated feedback on edge cases, high cognitive friction when switching between code editors and theory slides, and the financial cost and operational overhead of maintaining remote code execution sandboxes at scale.

---

## The Solution

This project provides a modern, gamified learning platform tailored to university syllabi, eliminating server-side execution risks and compute costs:

- **Client-side code execution**: Utilizes **Pyodide** (Python compiled to WebAssembly) to execute and test student solutions directly inside the browser, ensuring near-instant validation, zero backend compute costs, and isolation from host systems.
- **Curated challenge syllabus**: Multi-unit curriculum featuring progressive problem sets, guided starter templates, and automated unit test suites that evaluate edge cases in real time.
- **Built-in code editor**: Integrated with **Monaco Editor** (the engine powering VS Code) providing syntax highlighting, intelligent autocompletion, and familiar keyboard shortcuts.
- **Progress tracking & OAuth**: Secure authentication via **GitHub OAuth** and JWT sessions, recording completed challenges, test metrics, and student progress.

---

## Screenshots

| Home (logged out) | Home (logged in) |
|---|---|
| ![Home Unauth](images/apschool_home_unauth.png) | ![Home Auth](images/apschool_home_auth.png) |

| Challenges | Challenge solved |
|---|---|
| ![Challenges](images/apschool_challenges.png) | ![Test Passed](images/apschool_test_passed.png) |

| GitHub OAuth |
|---|
| ![GitHub OAuth](images/apschool_github_oauth.png) |

---

## Solution Architecture

The platform is deployed on AWS using a highly scalable, secure, and cost-effective cloud architecture provisioned entirely as code with Terraform. It combines a client-side WebAssembly execution runtime with a containerized Go REST API managed by an Auto Scaling Group behind an Application Load Balancer, distributed globally through Amazon CloudFront.

![AWS Cloud Architecture](diagrams/cloud_architecture.png)

### Architecture Components
- **Frontend SPA**: Single Page Application built with **Angular 21** and hosted on **Amazon S3**, distributed globally through **Amazon CloudFront** using Origin Access Control (OAC). Static assets are cached at edge locations, and all Python code execution occurs locally in the user's browser via WebAssembly.
- **Unified Reverse Proxy**: **Amazon CloudFront** provides a unified domain entry point, routing static requests (`/*`) to the S3 bucket and dynamic API traffic (`/api/*`) directly to the Application Load Balancer, eliminating Cross-Origin Resource Sharing (CORS) complexity.
- **Containerized REST API**: High-performance backend written in **Go 1.25** and packaged in **Docker** containers. Images are stored in **Amazon ECR** with automated lifecycle policies. The API runs on **EC2 instances** (Amazon Linux 2023) managed by an **Auto Scaling Group (ASG)** across public subnets, receiving traffic through an **Application Load Balancer (ALB)**.
- **Database Layer**: Managed relational database running **Amazon RDS PostgreSQL 16** in isolated private subnets across multiple availability zones. Ingress is strictly restricted on port 5432 to backend security groups.
- **Infrastructure as Code**: Modular and reproducible cloud infrastructure provisioned using **Terraform 1.16+**, utilizing remote state storage in **Amazon S3** with native lockfile state concurrency protection.
- **Monitoring & Observability**: Integrated with **Amazon CloudWatch** alarms tracking ALB unhealthy hosts, ASG CPU utilization, and target 5XX error rates, with automated notifications dispatched via **Amazon SNS**.

### CI/CD Automation Flow

The delivery process is fully automated via GitHub Actions, coordinating infrastructure changes, container builds, and frontend deployments:

![CI/CD Automation Flow](diagrams/cicd_flow.png)

- **Terraform Automation**: Provisions and manages the entire AWS infrastructure stack cleanly and idempotently.
- **Container Build & Registry**: Builds the Go backend binary into a hardened Docker image targeting `linux/amd64` and pushes to Amazon ECR.
- **Frontend Build & Edge Deployment**: Compiles the Angular 21 application via Bun, syncs assets to Amazon S3, and triggers an automated CloudFront cache invalidation.

---

## Repository Structure

```
apschool/
├── .github/workflows/    # CI test suite and deployment/teardown workflows
├── diagrams/             # Official AWS architecture and CI/CD diagrams
├── images/               # Application UI screenshots and assets
├── infra/                # Modular Infrastructure as Code (Terraform 1.16+)
│   ├── environments/prod # Production environment configuration and backend state
│   └── modules/          # VPC, RDS, ALB, Compute, ECR, Frontend, and Monitoring
├── server/               # Go 1.25 REST API, Goose database migrations, and challenge seeds
└── web/                  # Angular 21 SPA, Monaco Editor, and Pyodide WebAssembly runner
```

---

## Prerequisites
- Bun (for frontend build and dependencies)
- Go 1.25+ (for backend services)
- Docker (for container builds and local database)
- Terraform 1.16+ (for infrastructure management)

---

## Running Locally

1. **Database & Backend**:
   ```bash
   cd server
   docker compose up -d
   go run ./cmd/api
   ```

2. **Database Seeding**:
   ```bash
   cd server
   go run ./cmd/seed
   ```

3. **Frontend Client**:
   ```bash
   cd web
   bun install
   bun run start
   ```

4. Open `http://localhost:4200` in your browser.
