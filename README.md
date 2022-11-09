# script-go-aws-ecr-cleanup

Golang script to delete old untag ecr images

## .env configuration

``` bash
AWS_ECR_CLEANUP_SCRIPT_ID="yyyyyyyyyyyyyyyy"
AWS_ECR_CLEANUP_SCRIPT_KEY="xxxxxxxxxxxxxxxxxxxxxxxx"
AWS_ECR_CLEANUP_SCRIPT_REGION="us-east-1"
```

## run

go run main.go