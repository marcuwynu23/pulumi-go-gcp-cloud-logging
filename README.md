# Pulumi Go GCP Cloud Logging

This Pulumi project provisions Google Cloud Logging resources including a log bucket, log sink, and log-based metric.

## Resources Created

- **Log Bucket**: Cloud Logging storage bucket with configurable retention
- **Log Sink**: Routes logs matching a filter to the log bucket
- **Log-based Metric**: Creates a metric from logs matching a filter
- **API Enablement**: `logging.googleapis.com`

## Prerequisites

- Pulumi CLI installed
- Go 1.21+
- GCP credentials configured
- A GCP project with billing enabled

## Configuration

Copy `Pulumi.dev.yaml.example` to `Pulumi.dev.yaml` and update values:

```yaml
config:
  gcp:project: your-project-id
  gcp:region: us-central1
  logging:bucketName: default-log-bucket
  logging:retentionDays: "30"
  logging:sinkName: my-log-sink
  logging:sinkFilter: "severity >= ERROR"
  logging:metricName: error-count-metric
  logging:metricDescription: Count of ERROR severity logs
  logging:metricFilter: severity=ERROR
  logging:metricKind: DELTA
```

### Config Options

| Key | Description | Default |
|-----|-------------|---------|
| `gcp:project` | GCP Project ID | **Required** |
| `gcp:region` | GCP Region (Free Tier: us-central1, us-east1, us-west1) | `us-central1` |
| `logging:bucketName` | Log bucket ID | `default-log-bucket` |
| `logging:retentionDays` | Log retention in days | `30` |
| `logging:sinkName` | Log sink name | `my-log-sink` |
| `logging:sinkFilter` | Filter for log sink | `severity >= ERROR` |
| `logging:metricName` | Log-based metric name | `error-count-metric` |
| `logging:metricDescription` | Metric description | `Count of ERROR severity logs` |
| `logging:metricFilter` | Filter for metric | `severity=ERROR` |
| `logging:metricKind` | Metric kind (DELTA, GAUGE, CUMULATIVE) | `DELTA` |

## Deployment

```bash
# Install dependencies
go mod tidy

# Preview changes
pulumi preview

# Deploy
pulumi up
```

## Outputs

| Output | Description |
|--------|-------------|
| `logBucketId` | The ID of the log bucket |
| `logBucketLocation` | The location of the log bucket |
| `sinkName` | The name of the log sink |
| `metricName` | The name of the log-based metric |

## Free Tier Notes

Cloud Logging free tier includes:
- 50 GB logs ingested/month
- 30 days default retention (customizable)

See [Cloud Logging Pricing](https://cloud.google.com/logging/pricing) for details.

## Cleanup

```bash
pulumi destroy
```