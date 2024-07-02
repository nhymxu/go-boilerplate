# Prometheus metrics

In prometheus, you can set the `Authorization` header on every scrape request with the configured username and password.
Check the [scrape_configs](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config) configuration

```yaml
scrape_configs:
  - job_name: my_app
    basic_auth:
        username: admin!
        password: secret!
    static_configs:
        - targets:
            - 'localhost:8000'
```
