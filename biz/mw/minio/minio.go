package minio

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
	"offer_tiktok/pkg/constants"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	Client *minio.Client
	err    error
)

func MakeBucket(ctx context.Context, bucketName string) {
	exists, err := Client.BucketExists(ctx, bucketName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !exists {
		err = Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Successfully created mybucket %v\n", bucketName)
	}
}


func PutToBucket(ctx context.Context, bucketName string, file *multipart.FileHeader) (info minio.UploadInfo, err error) {
	fileObj, _ := file.Open()
	info, err = Client.PutObject(ctx, bucketName, file.Filename, fileObj, file.Size, minio.PutObjectOptions{})
	fileObj.Close()
	return info, err
}

func GetObjURL(ctx context.Context, bucketName string, filename string) (u *url.URL, err error) {
	exp := time.Hour * 24
	reqParams := make(url.Values)
	u, err = Client.PresignedGetObject(ctx, bucketName, filename, exp, reqParams)
	return u, err
}

func PutToBucketByBuf(ctx context.Context, bucketName string, filename string, buf *bytes.Buffer) (info minio.UploadInfo, err error) {
	info, err = Client.PutObject(ctx, bucketName, filename, buf, int64(buf.Len()), minio.PutObjectOptions{})
	return info, err
}

func Init() {
	ctx := context.Background()
	Client, err = minio.New(constants.MinioEndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(constants.MinioAccessKeyID, constants.MinioSecretAccessKey, ""),
		Secure: constants.MiniouseSSL})
	if err != nil {
		log.Fatalln("minio连接错误: ", err)
	}

	log.Printf("%#v\n", Client)

	MakeBucket(ctx, constants.MinioVideoBucketName)
	MakeBucket(ctx, constants.MinioImgBucketName)
}
