# Icbaat (I can be an artist too!)

Framework for creating non-generated web applications focused on productivity and fast development.
I head for personalized, opinionated, non-generated, deployable to [Zerops](https://zerops.io). Config will be taken from
environment variables.

## Tasks

### skelet (open source it)

- [X] commander -> cobra
- [X] configurator -> viper
- [X] di -> dig
- [X] app runner -> custom

### rest

- [X] http api server -> gin

### others

- [X] orm -> gorm
- [X] logging -> logrus
- [X] linting

### tests

- [ ] native go tests (calling straight the middlewares)

1. I hold only secret values and environment dependent values in configuration

## Ideas, next steps

- [ ] embed default component (with logger with service field inside) to all components
- [ ] description for config values
- [ ] docs -> swagger

## Tech stack

I want to get to the release ASAP. I don't think, that I can do that with HTML served from backend.

- Golang web application (HA)
    - PostgreSQL as database (everything will be written into it, except files) (HA)
    - Object Storage for file storing
