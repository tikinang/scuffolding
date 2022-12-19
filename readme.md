# Scuffolding

Framework for creating non-generated web applications focused on productivity and fast development.
I head for personalized, opinionated, **non-generated**.
Config will be taken from environment variables.

## Tasks

### skelet (open source it)

- [X] commander -> cobra
- [X] configurator -> viper
- [X] di -> dig
- [X] app runner -> custom

### rest

- [X] http api server -> gin
- [ ] complete auth module middlewares
- [ ] middlewares
  - [ ] cors
  - [ ] session

### others

- [X] orm -> gorm
- [X] logging -> logrus
- [X] linting

### tests

- [X] native go tests

1. I hold only secret values and environment dependent values in configuration

## Ideas, next steps

- [ ] embed default component (with logger with service field inside) to all components
- [ ] description for config values
- [ ] docs -> swagger
- [ ] object storage handler

# Deploy

- [ ] make Docker files
