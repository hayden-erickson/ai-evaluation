# Task Requirements

Data Types
- User
    - ID
    - Profile image url
    - Name
    - Time zone
    - Phone number
    - Created At Date & Time
- Habit
    - User ID (references a User)
    - ID
    - Name
    - Description
    - Created At Date & Time
- Log
    - Habit ID (references a Habit)
    - Created At Date & Time
    - Notes

Using the above data types can you create a REST API in golang that uses a sqlite database for local development. 

For each data type please create
    - SQL migration file to create a database table
    - A REST API endpoint that handles the POST, PUT, DELETE, and GET methods

Each user can only view their own habits

The API should have 
- jwt based authentication
- request input validation
- all errors are logged
- all errors returned as HTTP status codes with appropriate error messages
- separation of concerns with separate respository layer and service layer (modular architecture)
- dependency injection
- Robust, idiomatic, clean code 
- single responsibility interfaces
- clear comments on all functions and conditionals
- Strong security practices such as input validation, RBAC, secure headers
- Use standard library where possible (net/http, context, log).
- README.md file with instructions to run

# Instructions

Create an initial empty git commit with the message 'starting api generation task' before generating any files.

Any time you run a shell command please create a git commit before doing so with a short descriptions of changes.
- The idea is to create a log of any changes made an the thought process behind them or errors that were addressed.

Generate all the files needed to meet the above requirements.

Run `go run main.go`
- If there are any errors, fix them and create a git commit for the fix with a short descriptive comment
- Continue trying to run `go run main.go` and fixing errors until it succeeds

Once `go run main.go` succeeds with no errors the task is complete and no further action is required.
Please do not stop until this command succeeds. 

