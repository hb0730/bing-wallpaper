package upload

import (
	"context"
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
	Region       string
	protocol     string
	endpoint     string
	accessKey    string
	accessSecret string
}

func readEnv() (S3Info, error) {
	protocol, err := getEnv("protocol", "https")
	bucket, err := getEnv("bucket", "")
	region, err := getEnv("region", "Auto")
	endpoint, err := getEnv("endpoint", "")
	endpoint = strings.TrimLeft(endpoint, "https://")
	endpoint = strings.TrimLeft(endpoint, "http://")
	accessKey, err := getEnv("access_key", "")
	accessSecret, err := getEnv("access_key_secret", "")
	return S3Info{
		Bucket:       bucket,
		Region:       region,
		protocol:     protocol,
		endpoint:     endpoint,
		accessKey:    accessKey,
		accessSecret: accessSecret,
	}, err
}
func getEnv(key, value string) (string, error) {
	_value, s := os.LookupEnv(key)
	if s && _value != "" && len(_value) > 0 {
		return _value, nil
	}
	if value != "" && len(value) > 0 {
		return value, nil
	}
	return "", nil
}
func newS3Client(info S3Info) (*s3.Client, error) {
	// create a new S3 client
	return s3.New(
		s3.Options{
			Region: info.Region,
			Credentials: aws.NewCredentialsCache(
				credentials.NewStaticCredentialsProvider(info.accessKey, info.accessSecret, ""),
			),
			BaseEndpoint: aws.String(info.protocol + "://" + info.endpoint),
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
				filepath.Join("bing-wallpaper", time.Now().Format("2006-01"), resp.Filename),
			),
			Body: wallpaperRes.Body,
		},
	)
	return err
}
