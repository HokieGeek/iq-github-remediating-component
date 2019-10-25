#! /bin/bash

name=iq-github-remediating-component
roleArn=$(aws --profile admin iam get-role --role-name lambda_basic_execution |jq '.Role.Arn'|sed 's/\"//g')

GOARCH=amd64 GOOS=linux go build . || exit 1

zip ${name}.zip ${name}
aws --profile admin lambda get-function --function-name ${name} && {
    aws --profile admin \
        lambda update-function-code \
        --function-name ${name} \
        --zip-file fileb://${name}.zip \
        --publish
} || {
    aws --profile admin \
        lambda create-function \
        --function-name ${name} \
        --zip-file fileb://${name}.zip \
        --handler ${name} \
        --runtime go1.x \
        --timeout 120 \
        --role ${roleArn}
}

rm ${name}.zip
go clean