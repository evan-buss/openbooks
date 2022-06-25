# OpenBooks Service

## GRPC / REST API

Allow searching for books from a local cached and indexed copy.

## Implementations

### SQLITE / Postgres
- Store the books and their origins
- Use SQLITE until we experience performance issues.


### Searching

- SQLITE / Postgres Full Text Search
 - The simplest solution but may not work very well. SQLITE is not as good as POSTGRES for FTS and don't want to use external deps yet
- Bleve (Golang Lucene)
 - First choice
- ElasticSearch
 - Last resort if the other options aren't working out
 - Probably has the best functionality but requires additional services. Expensive to host on AWS