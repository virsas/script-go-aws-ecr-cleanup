package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print(".env file does not exist, will get the variables from the environment")
	}
}

func main() {
	session, err := createSession()
	if err != nil {
		panic(err)
	}
	svc := ecr.New(session)

	repositories, err := getRepositories(svc)
	if err != nil {
		panic(err)
	}

	for _, repo := range repositories {
		images, err := getImages(svc, repo.RepositoryName, "UNTAGGED")
		if err != nil {
			panic(err)
		}

		for _, img := range images {
			err := deleteImage(svc, repo.RepositoryName, img.ImageDigest)
			if err != nil {
				panic(err)
			}
		}
	}

}

func createSession() (*session.Session, error) {
	var awsID string = ""
	awsIDValue, awsIDPresent := os.LookupEnv("AWS_ECR_CLEANUP_SCRIPT_ID")
	if awsIDPresent {
		awsID = awsIDValue
	} else {
		return nil, errors.New("missing ENV Variable - AWS_ECR_CLEANUP_SCRIPT_ID")
	}

	var awsKey string = ""
	awsKeyValue, awsKeyPresent := os.LookupEnv("AWS_ECR_CLEANUP_SCRIPT_KEY")
	if awsKeyPresent {
		awsKey = awsKeyValue
	} else {
		return nil, errors.New("missing ENV Variable - AWS_ECR_CLEANUP_SCRIPT_KEY")
	}

	var region string = "eu-west-1"
	regionValue, regionPresent := os.LookupEnv("AWS_ECR_CLEANUP_SCRIPT_REGION")
	if regionPresent {
		region = regionValue
	} else {
		return nil, errors.New("missing ENV Variable - AWS_ECR_CLEANUP_SCRIPT_REGION")
	}

	session, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(awsID, awsKey, ""),
	})
	if err != nil {
		return nil, err
	}
	return session, nil
}

func getRepositories(svc *ecr.ECR) ([]*ecr.Repository, error) {
	opt := &ecr.DescribeRepositoriesInput{}
	repositories, err := svc.DescribeRepositories(opt)
	if err != nil {
		return nil, err
	}

	return repositories.Repositories, nil
}

func getImages(svc *ecr.ECR, repo *string, filter string) ([]*ecr.ImageIdentifier, error) {

	opt := &ecr.ListImagesInput{
		RepositoryName: repo,
		Filter: &ecr.ListImagesFilter{
			TagStatus: &filter,
		},
	}

	images, err := svc.ListImages(opt)
	if err != nil {
		return nil, err
	}

	return images.ImageIds, nil
}

func deleteImage(svc *ecr.ECR, repo *string, img *string) error {
	fmt.Println("#### Deleting image", *img, "in repo", *repo, "####")
	imageToDelete := &ecr.ImageIdentifier{ImageDigest: img}

	opt := &ecr.BatchDeleteImageInput{
		RepositoryName: repo,
		ImageIds:       []*ecr.ImageIdentifier{imageToDelete},
	}
	_, err := svc.BatchDeleteImage(opt)
	if err != nil {
		return err
	}

	fmt.Println("... Image deleted!")

	return nil
}
