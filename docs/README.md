# Parquet Factory Documentation

Parquet Factory is a program that can read from several data sources, aggregate the data received from them and generate a set of Parquet files with the aggregated data, storing them in a S3 bucket.

It is used to generate different data aggregations in the CCX Internal Data Pipeline, reading data from Kafka topics.

## Documentation

- **[Architecture](architecture.md)** - System architecture and data flow
- **[Configuration](config.md)** - Configuration options and environment variables
- **[CI/CD](ci.md)** - Continuous integration and code quality checks
- **[Deployment](deployment.md)** - Online and local deployment instructions
- **[Testing](testing.md)** - Unit tests, integration tests, and local testing
- **[Metrics](metrics.md)** - Prometheus metrics and monitoring

## Quick Start

See the main [README](../README.md) for quick start instructions on local development.

## Resources

- [Architecture diagram](resources/parquet-factory_hl.png)
