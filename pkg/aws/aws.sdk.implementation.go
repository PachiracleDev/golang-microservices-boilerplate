package implements

import (
	"context"
	"fmt"
	appConfig "moov/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

type AwsSdkImplementation struct {
	rekognitionClient *rekognition.Client
	s3Client          *s3.Client
	stsClient         *sts.Client
	config            *appConfig.Config
}

func NewSDKImplementation(conf *appConfig.Config) (*AwsSdkImplementation, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))

	if err != nil {
		return &AwsSdkImplementation{}, fmt.Errorf("error creating AWS session: %w", err)
	}

	return &AwsSdkImplementation{
		rekognitionClient: rekognition.NewFromConfig(cfg),
		s3Client:          s3.NewFromConfig(cfg),
		stsClient:         sts.NewFromConfig(cfg),
		config:            conf,
	}, nil
}

func (u AwsSdkImplementation) GetRekognitionToken() (*types.Credentials, error) {
	policy := `{
		"Statement": [
			{
				"Sid": "Stmt1RekognitionAccess",
				"Effect": "Allow",
				"Action": ["rekognition:*"],
				"Resource": "*"
			}
		]
	}`
	// Crear una solicitud para obtener un token federado
	command := &sts.GetFederationTokenInput{
		Name:            aws.String("RekognitionWebToken"),
		Policy:          aws.String(policy),
		DurationSeconds: aws.Int32(3600), // 1 hora
	}

	// Obtener el token federado
	result, err := u.stsClient.GetFederationToken(context.TODO(), command)
	if err != nil {
		return nil, err
	}

	return result.Credentials, nil
}

func (u AwsSdkImplementation) GetS3Token(key string) (*types.Credentials, error) {

	policy := fmt.Sprintf(`{
		"Statement": [
			{
				"Sid": "Stmt1S3UploadAssets",
				"Effect": "Allow",
				"Action": ["s3:PutObject"],
				"Resource": ["arn:aws:s3:::%s/%s/*"]
			}
		]
	}`, u.config.AWS.S3Bucket, key)

	input := &sts.GetFederationTokenInput{
		Name:            aws.String("S3UploadWebToken"),
		Policy:          aws.String(policy),
		DurationSeconds: aws.Int32(3600), // 1 hora
	}

	req, err := u.stsClient.GetFederationToken(context.Background(), input)

	if err != nil {
		return nil, fmt.Errorf("error getting S3 token: %w", err)
	}

	return req.Credentials, nil
}
