# Prompt
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

Using the above data types can you create a REST API in golang that uses a mysql database and will be deployed to AWS using kubernetes. 

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
- Strong security practices such as input validation, RBAC, HTTPS, secure headers, Kubernetes secrets
- Complete AWS kubernetes deployment setup with logging, monitoring, and scalable architecture
- Use standard library where possible (net/http, context, log).
- README.md file with instructions to run locally and to deploy to AWS 
