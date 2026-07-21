# Deployment

## Infraestructura

- **Cloud**: Google Cloud Platform (GCP)
- **Compute**: Cloud Run (serverless containers)
- **Database**: CloudSQL (PostgreSQL)
- **Red**: VPC nativa
- **Containers**: Docker

## Terraform (`infra/terraform/`)

```
infra/
  terraform/
    main.tf           → provider config, backend
    variables.tf      → variables de entorno
    outputs.tf        → outputs
    modules/
      cloud-run/      → Cloud Run service module
      cloud-sql/      → Cloud SQL module
      vpc/            → VPC network module
      iam/            → IAM roles and service accounts
    environments/
      dev/            → dev environment
      staging/        → staging environment (future)
      prod/           → production environment (future)
```

## CI/CD (GitHub Actions)

Workflows en `.github/workflows/`:

- `deploy.yml` — build + push a Container Registry + deploy a Cloud Run
- `test.yml` — tests unitarios y de integración
- `lint.yml` — linting (Go vet, ESLint, Prettier)

## Docker (`deploy/docker/`)

- `docker-compose.yml` — desarrollo local
  - `gateway` — Go API Gateway
  - `geospatial` — Geospatial Engine
  - `postgres` — PostgreSQL
  - `travel` — Next.js frontend (dev)

## Deploy Manual

```bash
# Frontend
cd frontend/travel && npm run build && systemctl restart cytaxi-travel

# Backend
systemctl restart cytaxi-geospatial
```

## Próximos Pasos (Blueprint)

- NATS JetStream cluster
- Service mesh (Istio o Linkerd)
- Redis cluster para caché distribuida
- CDN para assets estáticos
- Multi-región (Latam)
