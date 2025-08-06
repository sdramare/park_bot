# create AWS Lambda function
!/bin/bash
set -e
# Variables
FUNCTION_NAME="park_bot_function"
ZIP_FILE="park_bot.zip"
REGION="eu-north-1"  
ROLE_NAME="park_bot_lambda_role"

GOOS=linux GOARCH=amd64 go build -o bootstrap -tags lambda.norpc main.go 
rm $ZIP_FILE || true
zip $ZIP_FILE bootstrap

# if role exists, if not create it
if aws iam get-role --role-name $ROLE_NAME --region $REGION 2>/dev/null; then
  echo "Role $ROLE_NAME already exists."
else
  echo "Creating role $ROLE_NAME."
  aws iam create-role --role-name $ROLE_NAME \
  --assume-role-policy-document '{"Version": "2012-10-17", "Statement": [{"Effect": "Allow", "Principal": {"Service": "lambda.amazonaws.com"}, "Action": "sts:AssumeRole"}]}' \
  --region $REGION
    # Attach the AWSLambdaBasicExecutionRole policy to the role
    aws iam attach-role-policy --role-name $ROLE_NAME \
    --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole \
    --region $REGION

    aws iam attach-role-policy --role-name $ROLE_NAME \
    --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaRole \
    --region $REGION
    # Wait for the role to be created
    sleep 10   
fi
# create IAM role for Lambda with execution permissions
 

# Create Lambda function if it does not exist
if aws lambda get-function --function-name $FUNCTION_NAME --region $REGION 2>/dev/null; then
  echo "Lambda function $FUNCTION_NAME already exists."
else
  echo "Creating Lambda function $FUNCTION_NAME."
  aws lambda create-function --function-name $FUNCTION_NAME \
  --zip-file fileb://$ZIP_FILE \
  --handler bootstrap \
  --runtime 'provided.al2023' \
  --role arn:aws:iam::$(aws sts get-caller-identity --query Account --output text):role/$ROLE_NAME \
  --region $REGION
fi

# Add environment variables
aws lambda update-function-configuration --function-name $FUNCTION_NAME \
  --environment "Variables={TWILIO_ACCOUNT_SID=<SID>,TWILIO_AUTH_TOKEN=<TOKEN>}" \
  --region $REGION
# Add Labmda Url
aws lambda create-function-url-config --function-name $FUNCTION_NAME \
  --auth-type NONE \
  --region $REGION
# Output the function URL
aws lambda get-function-url-config --function-name $FUNCTION_NAME --region $REGION --query "FunctionUrl" --output text
echo "Lambda function URL: $(aws lambda get-function-url-config --function-name $FUNCTION_NAME --region $REGION --query "FunctionUrl" --output text)"
echo "Lambda function $FUNCTION_NAME created successfully in region $REGION."
rm $ZIP_FILE || true
rm bootstrap || true
