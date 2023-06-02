package upload

import (
	"context"
	"errors"
	"github.com/TimothyYe/bing-wallpaper"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// S3Info S3 OSS properties
type S3Info struct {
	Bucket       string
	Regin        string
	protocol     string
	endpoint     string
	accessKey    string
	accessSecret string
}

func readEnv() (S3Info, error) {
	protocol, err := getEnv("protocol", "https")
	bucket, err := getEnv("bucket", "")
	regin, err := getEnv("regin", "Auto")
	endpoint, err := getEnv("endpoint", "")
	endpoint = strings.TrimLeft(endpoint, "https://")
	endpoint = strings.TrimLeft(endpoint, "http://")
	accessKey, err := getEnv("accessKey", "")
	accessSecret, err := getEnv("accessKeySecret", "")
	return S3Info{
		Bucket:       bucket,
		Regin:        regin,
		protocol:     protocol,
		endpoint:     endpoint,
		accessKey:    accessKey,
		accessSecret: accessSecret,
	}, err
}
func getEnv(key, value string) (string, error) {
	v, s := os.LookupEnv(key)
	if s {
		return v, nil
	}
	if value != "" {
		return value, nil
	}
	return "", errors.New("env not found")
}
func newS3Client(info S3Info) (*s3.Client, error) {
	return s3.New(
		s3.Options{
			Region:           info.Regin,
			Credentials:      aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(info.accessKey, info.accessSecret, "")),
			EndpointResolver: s3.EndpointResolverFromURL(info.protocol + "://" + info.endpoint),
		},
	), nil
}

// S3OssClient S3 OSS upload client
type S3OssClient struct {
	s3Info   S3Info
	s3Client *s3.Client
}

// NewS3OssClient create a new S3OssClient
func NewS3OssClient() (*S3OssClient, error) {
	s3Info, err := readEnv()
	if err != nil {
		return nil, err
	}
	client, err := newS3Client(s3Info)
	if err != nil {
		return nil, err
	}
	return &S3OssClient{
		s3Info:   s3Info,
		s3Client: client,
	}, nil
}
func (client *S3OssClient) Upload(info WallpaperUploadInfo) error {
	resp, err := bing_wallpaper.Get(info.index, info.market, info.resolution)
	if err != nil {
		return err
	}
	wallpaperRes, err := http.Get(resp.URL)
	if err != nil {
		return err
	}
	_, err = client.s3Client.PutObject(
		context.TODO(),
		&s3.PutObjectInput{
			Bucket: &client.s3Info.Bucket,
			Key: aws.String(
				filepath.Join(time.Now().Format("2006-01"), resp.Filename),
			),
			Body: wallpaperRes.Body,
		},
	)
	return err
}
