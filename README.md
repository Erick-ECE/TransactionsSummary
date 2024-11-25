# Transactions Summary System

A system that processes CSV transaction files through AWS infrastructure (S3, Lambda, RDS) and automatically generates and emails summary reports to registered account holders.

## Overview

This project provides an automated pipeline for processing financial transactions:
1. Users upload CSV files containing transaction data to AWS S3 using pre-signed URLs
2. A Lambda function is triggered upon file upload to process the transactions
3. Transaction data is stored in a MySQL database (RDS)
4. Account and user information is also stored in the DB
5. Summary reports are automatically generated and emailed to account holders


## Prerequisites

- curl command-line tool
- For generating new URLs:
  - Go programming environment
  - AWS CLI configured with appropriate permissions
  - Access to AWS services (S3, Lambda, RDS)
  - 
**IMPORTANT NOTE: to change the email destination or add more accounts it must be Updated DB, please contact me to Do the update required**

## Usage

### Option 1: Use Pre-generated URL (Valid for 7 days )

You can directly upload your CSV file using this pre-generated URL:

```bash
curl --upload-file transactions.csv "https://transactions-demo.s3.us-east-2.amazonaws.com/transactions.csv?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAQXPZDEK6GV3NILUK%2F20241125%2Fus-east-2%2Fs3%2Faws4_request&X-Amz-Date=20241125T061704Z&X-Amz-Expires=604800&X-Amz-SignedHeaders=host&x-id=PutObject&X-Amz-Signature=0233eb612643838661b2cdff527cde95967d4ab1e4082b9bb77be9b2865ff7e9"
```

### Option 2: Generate New Pre-signed URL (Requires AWS Permissions)

If you need to generate a new URL, follow these steps:

1. Ensure you have AWS CLI configured with appropriate permissions
2. Run the Go script to generate a new pre-signed URL:
```bash
go run generate_presigned_url.go transactions.csv
```
3. Use the generated URL with curl:
```bash
curl --upload-file transactions.csv "<generated-URL>"
```

### CSV File Format

Your transaction file should follow this format:

```csv
Date,Transaction,AccountId
7/15,+60.5,1
7/28,-10.3,1
```

- Date: Transaction date (M/DD format)
- Transaction: Amount with sign (+ for credit, - for debit)
- AccountId: Account identifier

## System Flow

1. User uploads CSV file to S3 using either the pre-generated URL or a newly generated one
2. S3 upload triggers Lambda function
3. Lambda function:
   - Processes CSV file
   - Uploads transaction records in RDS - MySQL DB
   - Generates account summary
   - Sends email report to registered account email

## Output

After processing, the system automatically sends a summary email to the registered email address for each account, containing transaction summary and monthly breakdown.

## Database

```mermaid
erDiagram
    ACCOUNTS ||--o{ TRANSACTIONS : has
    ACCOUNTS {
        varchar(255) id PK
        float debit_balance
        float credit_balance
        varchar(255) email
    }
    TRANSACTIONS {
        varchar(255) id PK
        varchar(255) account_id FK
        float amount
        date transaction_date
        enum type
    }
```


