package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	bucket  = "build-artifacts"
	project = kingpin.Arg("project", "Project name").Required().String()
	step    = kingpin.Arg("step", "Step name to fetch artifact from").Required().String()
	env     = kingpin.Arg("env", "Environment name").Required().String()
)

func prefix(project, step, env string) string {
	return fmt.Sprintf("%s_%s/%s_%s_%s%s", env, project, env, step, project, step)
}

func main() {
	kingpin.Parse()
	s3svc := s3.New(session.New(), &aws.Config{
		Region: aws.String("eu-central-1"),
	})
	prefix := prefix(*project, *step, *env)
	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	resp, err := s3svc.ListObjects(params)
	if err != nil {
		log.Fatal("Bucket or keypath not found!")
	}
	var builds []int
	for _, value := range resp.Contents {
		parts := strings.Split(*value.Key, "/")
		id := parts[2]
		i, err := strconv.Atoi(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't parse string into int %s", id)
			continue
		}
		builds = append(builds, i)
	}

	sort.Ints(builds)
	latest_build := builds[len(builds)-1]
	key_path := fmt.Sprintf("%s/%d/%s", prefix, latest_build, "artifacts.txt")

	// Print informative data to stderr channel.
	fmt.Fprintf(os.Stderr, "Build ID: %d\n", latest_build)
	fmt.Fprintf(os.Stderr, "Artifacts Bucket: %s\n", bucket)
	fmt.Fprintf(os.Stderr, "Artifact Path: %s\n", key_path)

	result, err := s3svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key_path),
	})
	if err != nil {
		log.Fatal("Failed to get object", err)
	}

	fmt.Fprintf(os.Stderr, "Last Modfied: %s\n", result.LastModified)
	// Print only the artifact body to stdout channel.
	io.Copy(os.Stdout, result.Body)

	defer result.Body.Close()
}
