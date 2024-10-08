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
     "os/exec"
	 "strings"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

var (
	ctx           = context.Background()
	activeWorkers int
	mu            sync.Mutex
	maxWorkers    = 3
	s3Client      *s3.S3
	s3Uploader    *s3manager.Uploader
	s3Downloader  *s3manager.Downloader
	tempDir       = "./temp"
)

// Init initializes the worker package
func Init() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
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
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	s3Client = s3.New(awsSession)
	s3Uploader = s3manager.NewUploader(awsSession)
	s3Downloader = s3manager.NewDownloader(awsSession)

	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}

	log.Println("Worker package initialized successfully")
	return nil
}

func workerAvailable() bool {
	mu.Lock()
	defer mu.Unlock()
	if activeWorkers < maxWorkers {
		activeWorkers++
		return true
	}
	return false
}

func releaseWorker() {
	mu.Lock()
	activeWorkers--
	mu.Unlock()
}

func processVideo(s3Key string) {
	defer releaseWorker()

	downloadPath, err := downloadVideoFromS3(s3Key)
	if err != nil {
		log.Printf("Error downloading video %s: %v\n", s3Key, err)
		return
	}
	

	log.Printf("Processing video: %s\n", downloadPath)
	// Add your video processing logic here
	outputDir := filepath.Join(tempDir, "hls", filepath.Base(s3Key))
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Printf("Error creating output directory for %s: %v\n", s3Key, err)
		return
	}

	// Prepare FFmpeg command
	outputPath := filepath.Join(outputDir, "playlist.m3u8")
	cmd := exec.Command("ffmpeg",
		"-i", downloadPath,
		"-profile:v", "baseline", // baseline profile for wider device compatibility
		"-level", "3.0",
		"-start_number", "0",
		"-hls_time", "10", // 10 second segments
		"-hls_list_size", "0", // keep all segments
		"-f", "hls",
		outputPath,
	)

	// Run FFmpeg command
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error processing video %s: %v\nFFmpeg output: %s\n", s3Key, err, string(output))
		return
	}
	err = uploadHLSToS3(outputDir, s3Key)
	if err != nil {
		log.Printf("Error uploading HLS files for %s: %v\n", s3Key, err)
		return
	}

	
	log.Printf("Finished processing video: %s\n", s3Key)
}

func downloadVideoFromS3(s3Key string) (string, error) {
	log.Printf("Downloading video from S3: %s\n", s3Key)
	
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
func uploadHLSToS3(localDir, s3KeyPrefix string) error {
	err := filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %v", path, err)
		}
		defer file.Close()
		relativePath := strings.TrimPrefix(path, localDir)  
		s3Key := filepath.Join(s3KeyPrefix, relativePath)
		_, err = s3Uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(os.Getenv("S3_BUCKET_PUBLIC")),
			Key:    aws.String(s3Key),
			Body:   file,
		})
		if err != nil {
			return fmt.Errorf("failed to upload file %s: %v", path, err)
		}

		log.Printf("Uploaded %s to S3: %s\n", path, s3Key)
		return nil
	})

	return err
}
func Worker() {
	if err := Init(); err != nil {
		log.Fatalf("Failed to initialize worker: %v", err)
	}

	log.Println("Starting worker pool...")

	rdb := rediss.GetRedisClient()
	if rdb == nil {
		log.Fatal("Redis client not initialized")
	}

	for {
		job, err := rdb.BLPop(ctx, 5*time.Second, "test2").Result()
		if err != nil {
			
			continue
		}

		if len(job) < 2 {
			log.Println("Invalid job format received from Redis")
			continue
		}

		s3Key := job[1]
		log.Printf("Received job for video: %s\n", s3Key)

		for !workerAvailable() {
			log.Println("All workers busy, waiting for availability...")
			time.Sleep(1 * time.Second)
		}

		go processVideo(s3Key)
		log.Printf("Assigned job to worker for video: %s\n", s3Key)
	}
}