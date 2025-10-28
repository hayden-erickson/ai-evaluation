import os
import mysql.connector
from twilio.rest import Client
from datetime import datetime, timedelta
import logging

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

def get_db_connection():
    """Establishes a connection to the MySQL database."""
    try:
        connection = mysql.connector.connect(
            host=os.environ.get('MYSQL_HOST'),
            user=os.environ.get('MYSQL_USER'),
            password=os.environ.get('MYSQL_PASSWORD'),
            database=os.environ.get('MYSQL_DATABASE')
        )
        return connection
    except mysql.connector.Error as err:
        logging.error(f"Error connecting to database: {err}")
        raise

def send_twilio_notification(phone_number, message):
    """Sends a notification via Twilio."""
    try:
        account_sid = os.environ.get('TWILIO_ACCOUNT_SID')
        auth_token = os.environ.get('TWILIO_AUTH_TOKEN')
        twilio_phone_number = os.environ.get('TWILIO_PHONE_NUMBER')
        client = Client(account_sid, auth_token)

        client.messages.create(
            to=phone_number,
            from_=twilio_phone_number,
            body=message
        )
        logging.info(f"Successfully sent notification to {phone_number}")
    except Exception as e:
        logging.error(f"Failed to send Twilio notification to {phone_number}: {e}")

def get_users_to_notify():
    """Queries the database to find users who need a notification."""
    connection = get_db_connection()
    cursor = connection.cursor(dictionary=True)

    two_days_ago = datetime.now() - timedelta(days=2)
    yesterday = datetime.now() - timedelta(days=1)

    query = """
    SELECT u.id, u.phone_number, u.name
    FROM User u
    LEFT JOIN Habit h ON u.id = h.user_id
    LEFT JOIN Log l ON h.id = l.habit_id AND l.created_at >= %s
    GROUP BY u.id, u.phone_number, u.name
    HAVING COUNT(l.id) = 0 OR (COUNT(l.id) = 1 AND MAX(l.created_at) < %s);
    """

    try:
        cursor.execute(query, (two_days_ago.strftime('%Y-%m-%d %H:%M:%S'), yesterday.strftime('%Y-%m-%d %H:%M:%S')))
        users = cursor.fetchall()
        return users
    except mysql.connector.Error as err:
        logging.error(f"Error executing database query: {err}")
        return []
    finally:
        cursor.close()
        connection.close()

def main():
    """Main function to run the habit tracking notification job."""
    logging.info("Starting habit tracker cron job...")
    users_to_notify = get_users_to_notify()

    if not users_to_notify:
        logging.info("No users to notify.")
        return

    for user in users_to_notify:
        message = f"Hi {user['name']}, you have missed your habit logs. Keep up the good work!"
        send_twilio_notification(user['phone_number'], message)

    logging.info("Habit tracker cron job finished.")

if __name__ == "__main__":
    main()
