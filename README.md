# Resume Analytics
This repository holds source code for my metrics gathering for my resume website.

## Visitor Count
`count.go` is an AWS Lambda function that receives requests from API Gateway and interacts with DynamoDB.

- `GET /count`: Retrieve the current visitor count
- `PUT /count`: Increment the count by one
- `OPTIONS /count`: Support CORS
    - AWS has plenty of documentation and options to enable CORS support on API Gateway. It should add CORS headers to responses for you. This did not work for my situation :(
    - As such, we return the default HTTP headers that are expected in a CORS handshake
