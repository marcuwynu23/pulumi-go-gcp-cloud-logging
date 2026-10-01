package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/logging"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		project := config.Get(ctx, "gcp:project")
		if project == "" {
			return fmt.Errorf("gcp:project is required")
		}

		region := config.Get(ctx, "gcp:region")
		if region == "" {
			region = "us-central1"
		}

		bucketName := config.Get(ctx, "logging:bucketName")
		if bucketName == "" {
			bucketName = "default-log-bucket"
		}

		retentionDays := config.GetInt(ctx, "logging:retentionDays")
		if retentionDays == 0 {
			retentionDays = 30
		}

		sinkName := config.Get(ctx, "logging:sinkName")
		if sinkName == "" {
			sinkName = "my-log-sink"
		}

		sinkFilter := config.Get(ctx, "logging:sinkFilter")
		if sinkFilter == "" {
			sinkFilter = "severity >= ERROR"
		}

		metricName := config.Get(ctx, "logging:metricName")
		if metricName == "" {
			metricName = "error-count-metric"
		}

		metricDescription := config.Get(ctx, "logging:metricDescription")
		if metricDescription == "" {
			metricDescription = "Count of ERROR severity logs"
		}

		metricFilter := config.Get(ctx, "logging:metricFilter")
		if metricFilter == "" {
			metricFilter = "severity=ERROR"
		}

		metricKind := config.Get(ctx, "logging:metricKind")
		if metricKind == "" {
			metricKind = "DELTA"
		}

		// Enable Logging API
		_, err := projects.NewService(ctx, "logging-api", &projects.ServiceArgs{
			Service:          pulumi.String("logging.googleapis.com"),
			DisableOnDestroy: pulumi.Bool(false),
			Project:          pulumi.String(project),
		})
		if err != nil {
			return err
		}

		// Create log bucket
		logBucket, err := logging.NewProjectBucketConfig(ctx, "log-bucket", &logging.ProjectBucketConfigArgs{
			Project:         pulumi.String(project),
			Location:        pulumi.String(region),
			RetentionDays:   pulumi.Int(retentionDays),
			BucketId:        pulumi.String(bucketName),
			EnableAnalytics: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		// Create log sink
		sinkDestination := pulumi.Sprintf("logging.googleapis.com/projects/%s/buckets/%s", project, logBucket.BucketId)
		logSink, err := logging.NewProjectSink(ctx, "log-sink", &logging.ProjectSinkArgs{
			Name:                 pulumi.String(sinkName),
			Destination:          sinkDestination,
			Filter:               pulumi.String(sinkFilter),
			UniqueWriterIdentity: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		// Create log-based metric
		logMetric, err := logging.NewMetric(ctx, "log-metric", &logging.MetricArgs{
			Name:        pulumi.String(metricName),
			Description: pulumi.String(metricDescription),
			Filter:      pulumi.String(metricFilter),
			MetricDescriptor: &logging.MetricMetricDescriptorArgs{
				MetricKind: pulumi.String(metricKind),
				ValueType:  pulumi.String("INT64"),
			},
		})
		if err != nil {
			return err
		}

		// Export outputs
		ctx.Export("logBucketId", logBucket.BucketId)
		ctx.Export("logBucketLocation", logBucket.Location)
		ctx.Export("sinkName", logSink.Name)
		ctx.Export("metricName", logMetric.Name)

		return nil
	})
}
