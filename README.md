 # Pulumi serverless functions deployment in Go

## Overview
This little project demonstrates a simple (but robust) serverless deployment using Pulumi with AWS Lambda, API Gateway, and CloudWatch Logging.

## List of features
- Serverless Lambda function deployment
- Comprehensive error handling
- CloudWatch logging integration
- IAM role and policy management

## Prerequisites
- Go (1.16+)
- Pulumi CLI
- AWS Account, token
- AWS CLI configured

## Project Structure
```
.
├── main.go         # Primary Pulumi deployment script
├── function/       # Lambda function directory
│   └── handler.py  # Python Lambda handler
└── www/            # Static web content directory
```

## Configuration

### Lambda Function
- Runtime: Python 3.9
- Memory: 128 MB
- Timeout: 10 seconds
- Logging: CloudWatch integration

### API Gateway
- Routes:
  - `/`: Static content serving
  - `/date`: Lambda function endpoint

## Deployment Steps

1. Initialize Pulumi Project
```bash
pulumi new aws-go
```

2. Install Dependencies
```bash
go mod tidy
```

3. Configure AWS Credentials
```bash
aws configure
```

4. Deploy Infrastructure
```bash
pulumi up
```

## Error Handling
This deployment includes comprehensive error logging:
- Detailed error messages
- CloudWatch Log Group creation
- IAM permissions for logging

## Security Considerations
- Minimal IAM permissions
- Explicit role assumption policy
- Logging and tracing enabled

## Customization
Easily modify:
- Lambda function configuration
- API Gateway routes
- Log retention periods

## Monitoring
- CloudWatch Log Group configured
- Retention set to 14 days (adjustable)

## Troubleshooting
- Check Pulumi logs
- Verify AWS credentials
- Ensure function code is valid

## Contributing
1. Fork the repository
2. Create feature branch
3. Commit changes
4. Push to branch
5. Create Pull Request

## License
MIT License
