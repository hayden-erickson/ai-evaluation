Task:
Generate a production-ready Kubernetes CronJob that performs the following:

Queries a MySQL database for habit logs across all users over the past 2 days.
Sends a notification via Twilio to each user’s phone number if:

The user has no logs over the last 2 days.
The user has 1 log on day 1, but no log on day 2.
If both days have logs, no notification should be sent.



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

Evaluation Criteria:

Please explain your reasoning and design choices clearly.
Justify how your solution meets production standards.
Highlight any trade-offs or assumptions made.
Optionally suggest improvements or alternatives.