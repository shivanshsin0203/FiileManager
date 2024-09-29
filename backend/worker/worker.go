package worker

import (
	"context"
	"filemanager/redis"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
    "path/filepath"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var rdb *redis.Client 
var ctx = context.Background()
var activeWorkers int = 0 // Track the number of active workers
var mu sync.Mutex         // Mutex to control access to the worker counter
var maxWorkers = 3        // Limit the number of concurrent workers
var(
	s3Client    *s3.S3
	s3Uploader  *s3manager.Uploader
	s3Downloader *s3manager.Downloader
	tempDir        = "./temp"
)
func workerAvailable() bool {
	mu.Lock()
	defer mu.Unlock()
	if activeWorkers < maxWorkers {
		activeWorkers++ // Reserve a slot for the new worker
		return true
	}
	return false
}
func processVideo(s3Key string) {
	defer func() {
		mu.Lock()
		activeWorkers-- // Free up the slot when the worker is done
		mu.Unlock()
	}()
	// Download the video from S3
	downloadPath, err := downloadVideoFromS3(s3Key)
	if err != nil {
		log.Printf("Error downloading video %s: %v\n", s3Key, err)
		return
	}
	// defer os.Remove(downloadPath)
	fmt.Println("Processing video: ", downloadPath)
}
func downloadVideoFromS3(s3Key string) (string, error) {
	log.Println("Downloading video from S3:", s3Key)
	
	downloadPath := filepath.Join(tempDir, filepath.Base(s3Key))
	file, err := os.Create(downloadPath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	_, err = s3Downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
			Key:    aws.String(s3Key),
		})
	if err != nil {
		return "", fmt.Errorf("error downloading file: %v", err)
	}

	return downloadPath, nil
}
func Worker(){
      rdb := rediss.GetRedisClient()
	  err := godotenv.Load()
	  if err != nil {
		  log.Fatal("Error loading .env file")
	  }
	  awsSession, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(os.Getenv("AWS_REGION")),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("AWS_ACCESS_KEY_ID"),
				os.Getenv("AWS_SECRET_ACCESS_KEY"),
				"",
			),
		},
	})

	if err != nil {
		panic(err)
	}
     
	s3Uploader = s3manager.NewUploader(awsSession)
	s3Downloader = s3manager.NewDownloader(awsSession)
	os.MkdirAll(tempDir, os.ModePerm)
	  for{
	  job, err := rdb.BLPop(ctx, 0*time.Second, "test2").Result()
		if err != nil {
			log.Println("Error fetching from Redis:", err)
			continue
		}
			// Check if a worker is available
			if workerAvailable() {
				s3Key := job[1] 
				go processVideo(s3Key) 
				fmt.Println("Assigned job to worker for video:", s3Key)
			} else {
				log.Println("All workers busy, waiting for availability...")
				time.Sleep(1 * time.Second) // Wait and retry if no workers are available
			}
	}
}