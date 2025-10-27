Task:
Generate a production-ready Kubernetes CronJob that performs the following:

Queries a MySQL database for habit logs across all users over the past 2 days.
Sends a notification via Twilio to each user’s phone number if:

The user has no logs over the last 2 days.
The user has 1 log on day 1, but no log on day 2.
If both days have logs, no notification should be sent.

These are the tables the database has
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

Requirements:

Include all necessary Kubernetes configuration files:

CronJob manifest
Secrets/ConfigMaps
RBAC (if applicable)
Resource limits and probes


Use secure practices for handling credentials and secrets.
Ensure the solution is deployable to a real Kubernetes cluster.
Include logging and error handling.
Make the code modular and maintainable.