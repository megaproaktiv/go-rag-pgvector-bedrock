# API for RAG with embedding

![overview](img/architecture.png)

## Manually create RDS/postgres

Minimal Version 15.3

## Import

See directory `import`.


## Deploy API

```bash
cd api
task build 
task deploy
```

## Test

In Postman:

- Get URL from SAM output
- Put API Gateway API Key as autorization header X-API-Key
- Method POST

Body raw:

```json
{
    "question": "This is my question?"
}
```

## Include in Webserver