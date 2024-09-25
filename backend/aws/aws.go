package aws

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "time"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/joho/godotenv"
)

const urlExpiration = 15 * time.Minute

func init() {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        fmt.Println("Error loading .env file")
    }
}

func GeneratePresignedURL(w http.ResponseWriter, r *http.Request) {
    // Get credentials and region from environment variables
    accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
    secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
    region := os.Getenv("AWS_REGION")
    bucketName := os.Getenv("S3_BUCKET_NAME")
    fmt.Println(accessKeyID, secretAccessKey, region, bucketName)
    if accessKeyID == "" || secretAccessKey == "" || region == "" || bucketName == "" {
        http.Error(w, "Missing AWS credentials or S3 bucket", http.StatusInternalServerError)
        return
    }

    // Manually configure AWS with credentials from environment variables
    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithRegion(region),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
    )
    if err != nil {
        http.Error(w, "Failed to load AWS config: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Create an S3 client
    s3Client := s3.NewFromConfig(cfg)

    // Get file name from request
    fileName := r.URL.Query().Get("fileName")
    if fileName == "" {
        http.Error(w, "fileName is required", http.StatusBadRequest)
        return
    }

    // Create a presigned URL for uploading
    presignClient := s3.NewPresignClient(s3Client)
    req, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(fileName),
    }, s3.WithPresignExpires(urlExpiration)) // URL expiration time
    if err != nil {
        http.Error(w, "Failed to generate pre-signed URL: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Return the presigned URL to the client
    fmt.Fprintf(w, req.URL)
}
