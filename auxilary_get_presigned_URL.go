package main

// import (
// 	"context"
// 	"time"

// 	"github.com/aws/aws-sdk-go-v2/service/s3"
// )

// func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {

// 	psS3Client := s3.NewPresignClient(s3Client)

// 	psGettedObj, err := psS3Client.PresignGetObject(context.TODO(), &s3.GetObjectInput{Bucket: &bucket, Key: &key}, s3.WithPresignExpires(expireTime))
// 	if err != nil {
// 		return "", err
// 	}
// 	return psGettedObj.URL, nil
// }
