!/bin/bash
set -e
#   
FUNCTION_NAME="park_bot_function"
ZIP_FILE="park_bot.zip"
REGION="eu-north-1"

GOOS=linux GOARCH=amd64 go build -o bootstrap -tags lambda.norpc main.go 
rm $ZIP_FILE || true
zip $ZIP_FILE bootstrap

aws lambda update-function-code \
    --function-name $FUNCTION_NAME \
    --zip-file fileb://$ZIP_FILE
rm $ZIP_FILE || true
rm bootstrap || true