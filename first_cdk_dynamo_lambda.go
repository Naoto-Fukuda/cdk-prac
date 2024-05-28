package main

import (
	"log"
	"os"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/joho/godotenv"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
)

type AppStackProps struct {
	awscdk.StackProps
}

func NewAppStack(scope constructs.Construct, id string, props *AppStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)

	// DynamoDBテーブルの作成
	table := awsdynamodb.NewTable(stack, jsii.String("MyDynamoDB"), &awsdynamodb.TableProps{
		TableName: jsii.String("MyDynamoDB"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode: awsdynamodb.BillingMode_PAY_PER_REQUEST,
	})

	// Lambda関数の定義
	lambdaFunction := awslambda.NewFunction(stack, jsii.String("MyFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2(),
		Handler: jsii.String("main"),
		Code:    awslambda.Code_FromAsset(jsii.String("lambda/dynamoDBHandler"), nil),
		Environment: &map[string]*string{
			"DYNAMODB_TABLE_NAME": table.TableName(),
		},
	})

	// テーブルへの権限付与
	table.GrantReadWriteData(lambdaFunction)

	// API Gatewayの定義
	api := awsapigatewayv2.NewHttpApi(stack, jsii.String("MyApi"), &awsapigatewayv2.HttpApiProps{
		ApiName: jsii.String("MyApi"),
	})

	// Lambda関数の統合
	lambdaIntegration := awsapigatewayv2integrations.NewHttpLambdaIntegration(jsii.String("MyIntegration"), lambdaFunction, nil)
	
	// API Gatewayの定義
	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path: jsii.String("/api-endpoint"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_GET,
			awsapigatewayv2.HttpMethod_POST,
			awsapigatewayv2.HttpMethod_PUT,
			awsapigatewayv2.HttpMethod_DELETE,
		},
		Integration: lambdaIntegration,
	})

	return stack
}

func main() {
	// .envファイルを読み込む
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	// スタックの入れ物を作るイメージ
	app := awscdk.NewApp(nil)

	NewAppStack(app, "AppStack", &AppStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	awsAccountId := os.Getenv("AWS_ACCOUNT_ID")
	awsRegion := os.Getenv("AWS_REGION")

	return &awscdk.Environment{
		Account: jsii.String(awsAccountId),
		Region:  jsii.String(awsRegion),
	}
}
