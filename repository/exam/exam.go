package exam

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetExamFiles(ctx context.Context, examID int) ([]string, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("eu-west-1"),
	)
	if err != nil {
		return nil, err
	}

	if err != nil {
		log.Fatalf("failed to create AWS session, %v", err)
	}

	awsS3Client := s3.NewFromConfig(cfg)

	objects, err := awsS3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String("datadon-data"),
		Prefix: aws.String(fmt.Sprintf("exam/%d/", examID)),
	})
	if err != nil {
		return nil, err
	}

	log.Printf("%v", objects)

	list := make([]string, 0)
	for _, obj := range objects.Contents {
		key := *obj.Key
		key = strings.Replace(key, fmt.Sprintf("exam/%d/", examID), "", 1)
		if key != "" {
			list = append(list, key)
		}
	}

	return list, nil
}
