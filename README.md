Take Home Assignment - Build a tiny ledger! 
===========================================

Simple API for a simple ledger that allows users to do deposits, withdrawals, view current balance and view transaction 
history.

The application makes use of an in-memory database to store accounts and transactions (deposits and withdrawals).

For simplicity, account creation is done on account retrieval or on a valid deposit. 
Calling list transactions or making a withdrawal with an account id that doesn't exist results in a not found error 
being returned.

How to run
==========

With go installed cd to the root of the repository and run `go run .` to start the application. 
You will need to have port 8080 free to run the application.

Using your preferred api tool, you can call the various api endpoints. The endpoints along with example payloads and 
response are described below.

Example requests
================

Get account/balance
```http request
### Expect {"account":{"ID":"0193b58c-4bab-7c00-924f-c51f27ec2e05","Balance":0}}
GET http://localhost:8080/account/0193b58c-4bab-7c00-924f-c51f27ec2e05
```

Make a deposit
```http request
### Expect {"transaction":{"TransactionID":"2792f482-7ea5-44ab-8f7b-725303070326","Amount":20.5,"TransactionTime":"2025-02-09T17:51:53.917663Z","TransactionType":"deposit"}}
POST http://localhost:8080/account/0193b58c-4bab-7c00-924f-c51f27ec2e05/deposit
Content-Type: application/json

   {
    "amount": 20.50
    }
```

Make a withdrawal
```http request
### Expect {"transaction":{"TransactionID":"2792f482-7ea5-44ab-8f7b-725303070326","Amount":20.5,"TransactionTime":"2025-02-09T17:51:53.917663Z","TransactionType":"withdrawal"}}
POST http://localhost:8080/account/0193b58c-4bab-7c00-924f-c51f27ec2e05/withdrawal
Content-Type: application/json

   {
    "amount": 20.50
    }
```

List transactions (deposits and withdrawals)
```http request
### Expect {"next_page":0,"transactions":[{"TransactionID":"9771738b-c582-4b0a-9cc6-88cc0fec2b33","Amount":20.5,"TransactionTime":"2025-02-09T18:29:03.812856Z","TransactionType":"deposit"},{"TransactionID":"98db0e52-8fe7-4aae-946a-f95aade54893","Amount":20.5,"TransactionTime":"2025-02-09T18:29:02.753119Z","TransactionType":"deposit"},{"TransactionID":"a8e9c58d-1182-4a8b-9d8f-ce47162a21f8","Amount":20.5,"TransactionTime":"2025-02-09T18:29:01.822012Z","TransactionType":"deposit"}]}
GET http://localhost:8080/account/0193b58c-4bab-7c00-924f-c51f27ec2e05/transactions
Content-Type: application/json

   {
    "limit": 10,
    "page": 0
    }
```