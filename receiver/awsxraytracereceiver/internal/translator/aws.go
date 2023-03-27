package translator

import (
	awsxray "github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/xray"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"
)

type AWSBuilder struct {
	OTLPSpanResourceBuilder
}

func (a *AWSBuilder) WithAws(data *awsxray.AWSData) *AWSBuilder {
	if data != nil {
		if data.EC2 != nil {
			a.pResource.Attributes().PutStr(conventions.AttributeHostID, *data.EC2.InstanceID)
			a.pResource.Attributes().PutStr(conventions.AttributeCloudAvailabilityZone, *data.EC2.AvailabilityZone)
			a.pResource.Attributes().PutStr(conventions.AttributeHostType, *data.EC2.InstanceSize)
			a.pResource.Attributes().PutStr(conventions.AttributeHostImageID, *data.EC2.AmiID)
		}
		if data.ECS != nil {
			a.pResource.Attributes().PutStr(conventions.AttributeContainerName, *data.ECS.ContainerName)
			a.pResource.Attributes().PutStr(conventions.AttributeContainerID, *data.ECS.ContainerID)
			a.pResource.Attributes().PutStr(conventions.AttributeCloudAvailabilityZone, *data.ECS.AvailabilityZone)
			a.pResource.Attributes().PutStr(conventions.AttributeAWSECSContainerARN, *data.ECS.ContainerArn)
			a.pResource.Attributes().PutStr(conventions.AttributeAWSECSClusterARN, *data.ECS.ClusterArn)
			a.pResource.Attributes().PutStr(conventions.AttributeAWSECSTaskARN, *data.ECS.TaskArn)
			a.pResource.Attributes().PutStr(conventions.AttributeAWSECSTaskFamily, *data.ECS.TaskFamily)
			a.pResource.Attributes().PutStr(conventions.AttributeAWSECSLaunchtype, *data.ECS.LaunchType)
		}
		if data.Beanstalk != nil {
			a.pResource.Attributes().PutStr(conventions.AttributeServiceNamespace, *data.Beanstalk.Environment)
			a.pResource.Attributes().PutInt(conventions.AttributeServiceInstanceID, *data.Beanstalk.DeploymentID)
			a.pResource.Attributes().PutStr(conventions.AttributeServiceVersion, *data.Beanstalk.VersionLabel)
		}
		if data.EKS != nil {
			a.pResource.Attributes().PutStr(conventions.AttributeK8SClusterName, *data.EKS.ClusterName)
			a.pResource.Attributes().PutStr(conventions.AttributeK8SPodName, *data.EKS.Pod)
			a.pResource.Attributes().PutStr(conventions.AttributeContainerID, *data.EKS.ContainerID)
		}
		if data.XRay != nil {
			if data.XRay.SDK != nil {
				a.pResource.Attributes().PutStr(conventions.AttributeTelemetrySDKName, *data.XRay.SDK)
			}
			if data.XRay.SDKVersion != nil {
				a.pResource.Attributes().PutStr(conventions.AttributeTelemetrySDKVersion, *data.XRay.SDKVersion)
			}
			if data.XRay.AutoInstrumentation != nil {
				a.pResource.Attributes().PutBool(conventions.AttributeTelemetryAutoVersion, *data.XRay.AutoInstrumentation)
			}
		}
		if data.Operation != nil {
			a.pResource.Attributes().PutStr(awsxray.AWSOperationAttribute, *data.Operation)
		}
		if data.RemoteRegion != nil {
			a.pResource.Attributes().PutStr(awsxray.AWSRegionAttribute, *data.RemoteRegion)
		}
		if data.RequestID != nil {
			a.pResource.Attributes().PutStr(awsxray.AWSRequestIDAttribute2, *data.RequestID)
		}
		if data.QueueURL != nil {
			a.pResource.Attributes().PutStr(awsxray.AWSQueueURLAttribute2, *data.QueueURL)
		}
		if data.TableName != nil {
			a.pResource.Attributes().PutStr(awsxray.AWSTableNameAttribute2, *data.TableName)
		}
	}
	return a
}
