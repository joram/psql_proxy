# psql anonymizing proxy
This is a proxy for a PostgreSQL database.
There will be an anonymization layer to remove PII from any select statements.
## Intended Use Case
this is intended to be used for production databases that have PII in them, where you wish to analyize the data, but do not want to expose the PII to the analysts.

## Getting Started
### Prerequisites
* docker
* go 1.20+
* make
* psql

### To Run
1. `docker compose up -d`
2. `make run`

### To Demo the SQL example files
1. `make test_queries`

