# scaling-adtech

In this repository there is an example of a simple ad server in golang.

Starts it with
```bash
docker compose up
```

Then you can check the health of the service with
```bash
curl http://localhost:8080/health
```

See the [OpenAPI document](api/openapi.yaml) for the complete API specification.


# Design decision

- Postgres with prepared statements
- PubSub
- Low tech


# Scaling considerations

- Increase GO CPUs
- Scale up postgres
- Improve postgres
- Use Redis for caching
- Push architecture, RocksDB

# TODO

- [x] Postgres setup
- [x] Postgres connection
- [ ] Validation
- [x] Winning ad
- [x] Tracking