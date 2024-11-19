package main

import (
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi-aws-apigateway/sdk/v2/go/apigateway"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createLambdaLogGroup(ctx *pulumi.Context, functionName string) (*cloudwatch.LogGroup, error) {
	return cloudwatch.NewLogGroup(ctx, "lambda-log-group", &cloudwatch.LogGroupArgs{
		Name:            pulumi.String(functionName),
		RetentionInDays: pulumi.Int(14),
	})
}

func createLambdaExecutionRole(ctx *pulumi.Context, logGroup *cloudwatch.LogGroup) (*iam.Role, error) {
	assumeRolePolicy, err := json.Marshal(map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Action": "sts:AssumeRole",
				"Effect": "Allow",
				"Principal": map[string]interface{}{
					"Service": "lambda.amazonaws.com",
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// add your inline policy for CloudWatch Logs
	logPolicy, err := iam.NewRolePolicy(ctx, "lambda-log-policy", &iam.RolePolicyArgs{
		Role: pulumi.String(fmt.Sprintf("%s-log-policy", logGroup.Name)),
		Policy: pulumi.String(fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Action": [
						"logs:CreateLogStream",
						"logs:PutLogEvents"
					],
					"Resource": "%s:*"
				}
			]
		}`, logGroup.Arn)),
	})
	if err != nil {
		return nil, err
	}

	role, err := iam.NewRole(ctx, "lambda-execution-role", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(assumeRolePolicy),
		ManagedPolicyArns: pulumi.StringArray{
			iam.ManagedPolicyAWSLambdaBasicExecutionRole,
		},
	})
	if err != nil {
		return nil, err
	}

	return role, nil
}

func createLambdaFunction(ctx *pulumi.Context, role *iam.Role, logGroup *cloudwatch.LogGroup) (*lambda.Function, error) {
	fn, err := lambda.NewFunction(ctx, "date-function", &lambda.FunctionArgs{
		Runtime:    pulumi.String("python3.9"),
		Handler:    pulumi.String("handler.handler"),
		Role:       role.Arn,
		Code:       pulumi.NewFileArchive("./function"),
		MemorySize: pulumi.Int(128),
		Timeout:    pulumi.Int(10),
	})
	if err != nil {
		return nil, err
	}

	return fn, nil
}

func createRestAPI(ctx *pulumi.Context, fn *lambda.Function) (*apigateway.RestAPI, error) {
	localPath := "www"
	method := apigateway.MethodGET

	api, err := apigateway.NewRestAPI(ctx, "date-api", &apigateway.RestAPIArgs{
		Routes: []apigateway.RouteArgs{
			{
				Path:      "/",
				LocalPath: &localPath,
			},
			{
				Path:         "/date",
				Method:       &method,
				EventHandler: fn,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return api, nil
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create CloudWatch Log Group
		logGroup, err := createLambdaLogGroup(ctx, "/aws/lambda/date-function")
		if err != nil {
			return fmt.Errorf("failed to create log group: %w", err)
		}

		// Create Lambda execution role with logging permissions
		role, err := createLambdaExecutionRole(ctx, logGroup)
		if err != nil {
			return fmt.Errorf("failed to create IAM role: %w", err)
		}

		// Create Lambda function
		fn, err := createLambdaFunction(ctx, role, logGroup)
		if err != nil {
			return fmt.Errorf("failed to create Lambda function: %w", err)
		}

		// Create REST API
		api, err := createRestAPI(ctx, fn)
		if err != nil {
			return fmt.Errorf("failed to create API Gateway: %w", err)
		}

		// Export the API URL
		ctx.Export("url", api.Url)
		return nil
	})
}
