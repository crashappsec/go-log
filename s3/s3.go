package s3

// useful examples https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/gov2

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	AWS_ENDPOINT_URL = os.Getenv("AWS_ENDPOINT_URL")
	client           *s3Client
)

type s3Client struct {
	client    *s3.Client
	presigner *s3.PresignClient
}

type GetObjectInput = s3.GetObjectInput
type PutObjectInput = s3.PutObjectInput
type HeadObjectInput = s3.HeadObjectInput

type PresignOptions struct {
	ExpiresInSeconds int64
}

type ObjectExistsInput struct {
	Bucket string
	Key    string
}

func GetOrNew(ctx context.Context) *s3Client {
	if client == nil {
		var err error
		client, err = New(ctx)
		if err != nil {
			log.Fatalf("failed to init s3 client %s", err.Error())
		}
	}
	return client
}

func GetAWSConfig(ctx context.Context) (*aws.Config, error) {
	awsRegion := os.Getenv("AWS_DEFAULT_REGION")

	if awsRegion == "" {
		awsRegion = os.Getenv("AWS_REGION")
	}

	if awsRegion == "" {
		return nil, fmt.Errorf("AWS_REGION or AWS_DEFAULT_REGION must be set")
	}

	var cfg aws.Config
	var err error
	if AWS_ENDPOINT_URL != "" {
		resolver := aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           AWS_ENDPOINT_URL,
					SigningRegion: awsRegion,
				}, nil
			},
		)
		cfg, err = config.LoadDefaultConfig(
			ctx,
			config.WithEndpointResolverWithOptions(resolver),
		)

		if err != nil {
			return nil, err
		}
	} else {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))

		if err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}

func New(ctx context.Context) (*s3Client, error) {
	cfg, err := GetAWSConfig(ctx)
	if err != nil {
		return nil, err
	}
	var client *s3.Client
	if AWS_ENDPOINT_URL == "" {
		client = s3.NewFromConfig(*cfg)
	} else {
		client = s3.NewFromConfig(*cfg, func(o *s3.Options) {
			o.UsePathStyle = true
		})
	}
	presigner := s3.NewPresignClient(client)
	return &s3Client{
		client:    client,
		presigner: presigner,
	}, nil
}

func (c *s3Client) PutObject(
	ctx context.Context,
	input *PutObjectInput,
) (*s3.PutObjectOutput, error) {
	return c.client.PutObject(ctx, input)
}

func (c *s3Client) GetObject(
	ctx context.Context,
	input *GetObjectInput,
) (*s3.GetObjectOutput, error) {
	return c.client.GetObject(ctx, input)
}

func (c *s3Client) GetObjectBytes(
	ctx context.Context,
	input *GetObjectInput,
) ([]byte, error) {
	result, err := c.GetObject(ctx, input)
	if err != nil {
		return nil, err
	}

	defer result.Body.Close()

	bodyInBytes, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	return bodyInBytes, nil
}

func (c *s3Client) PresignGetObject(
	ctx context.Context,
	input *GetObjectInput,
	options *PresignOptions,
) (*v4.PresignedHTTPRequest, error) {
	request, err := c.presigner.PresignGetObject(
		ctx,
		input,
		func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(
				options.ExpiresInSeconds * int64(time.Second),
			)
		},
	)
	// if err != nil {
	// 	log.Errorf(
	// 		"Couldn't get a presigned request for GetObject",
	// 		log.String("bucket", *input.Bucket),
	// 		log.String("key", *input.Key),
	// 		log.Err(err),
	// 	)
	// }
	return request, err
}

func (c *s3Client) PresignPutObject(
	ctx context.Context,
	input *PutObjectInput,
	options *PresignOptions,
) (*v4.PresignedHTTPRequest, error) {
	request, err := c.presigner.PresignPutObject(
		ctx,
		input,
		func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(
				options.ExpiresInSeconds * int64(time.Second),
			)
		},
	)
	// if err != nil {
	// 	c.logger.Error(
	// 		"Couldn't get a presigned request for PutObject",
	// 		log.String("bucket", *input.Bucket),
	// 		log.String("key", *input.Key),
	// 		log.Err(err),
	// 	)
	// }
	return request, err
}

func (c *s3Client) ObjectExists(
	ctx context.Context,
	input *HeadObjectInput,
) (bool, error) {
	_, err := c.client.HeadObject(ctx, input)
	if err != nil {
		// https://aws.github.io/aws-sdk-go-v2/docs/handling-errors/
		// https://stackoverflow.com/questions/57697095/how-i-can-safely-check-if-file-exists-in-s3-bucket-using-go-in-lambda
		var responseError *awshttp.ResponseError
		if errors.As(err, &responseError) &&
			responseError.ResponseError.HTTPStatusCode() == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c *s3Client) PutObjectIfMissing(
	ctx context.Context,
	input *PutObjectInput,
) (bool, *s3.PutObjectOutput, error) {
	exists, err := c.ObjectExists(ctx, &HeadObjectInput{
		Bucket: input.Bucket,
		Key:    input.Key,
	})
	if err != nil {
		return false, nil, err
	}
	if exists {
		return false, nil, nil
	}
	output, err := c.PutObject(ctx, input)
	return true, output, err
}
