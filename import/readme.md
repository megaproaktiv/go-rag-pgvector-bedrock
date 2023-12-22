# Import Data with embeddings

- Filename hardcoded "testdata/testdata.pdf
- Get embeddings with bedrock/titan
- Needs an (Aurora) pgsql database 15.3


## Sentence-window retrieval

A text chunk is vectorized. Also stored is the context of the chunk. The context is the chunk before and after the chunk.

## Usage

Set environment:

# https://www.postgresql.org/docs/current/libpq-envars.html

```bash
export PGUSER="postgres"
export PGPASSWORD="***"
export PGHOST= "database-embedding.xyz.eu-west-1.rds.amazonaws.com"
export PGPORT= 5432
export PGDATABASE= "vectordb"
export PGCONNECT_TIMEOUT= 5
export AWS_REGION= "eu-central-1"
export AWS_DEFAULT_REGION= "eu-central-1"
```

Then run:

```bash
go run main.go  -d=true -f testdata/sample.pdf
```

Parameter | Description
--- | ---
-d | Drop tables before import
-f | Filename of pdf to import
