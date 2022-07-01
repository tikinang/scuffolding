# TIKIGO

Framework for creating non-generated REST APIs focused on productivity and fast development.
I head for personalized, opinionated, non-generated.
I head for deploy on [Zerops](https://zerops.io). Config will be taken from environment variables.
I head for skeleton, that can handle possible future non-API services and tools. It should have working command line interface.
Maybe in the future, I will open source (some of) it.

## TODO

### skelet (open source it)

* commander -> cobra ✅
* configurator -> viper ✅
* di -> dig ✅
* app runner -> custom ✅

### rest

* http api server -> gin ✅

I am going to need:

1. html serving + cash
2. api endpoints (for cash ajax calls)

### others

* orm -> gorm
* logging -> logrus ✅
* linting

### tests

* native go tests (calling straight the middlewares)

## Rules

1. I hold only secret values and environment dependent values in configuration

## Ideas, next steps

1. embed default component (with logger with service field inside) to all components
2. description for config values
3. docs -> swagger


## Tech stack

I want to get to the release ASAP. I don't think, that I can do that with HTML served from backend.

- Golang API (HA)
  - MariaDB as Database (everything will be written into it, except images) (HA)
- Svelte Web Application (HA)
  - pre-rendered / server-side-rendered entrypoints
  - Bulma as CSS Framework
- Nginx as Application Balancer and for HTTPS termination (HA) 
