#!/usr/bin/env python3
"""
Habit Tracking Notification Service
Queries MySQL database for habit logs and sends Twilio notifications to users
who haven't logged their habits in the past 2 days.
"""

import os
import sys
import logging
from datetime import datetime, timedelta
from typing import List, Dict, Optional
import pymysql
from twilio.rest import Client
from twilio.base.exceptions import TwilioRestException
import pytz

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.StreamHandler(sys.stdout)
    ]
)
logger = logging.getLogger(__name__)


class DatabaseConfig:
    """Database configuration from environment variables"""
    def __init__(self):
        self.host = os.getenv('DB_HOST', 'localhost')
        self.port = int(os.getenv('DB_PORT', '3306'))
        self.user = os.getenv('DB_USER')
        self.password = os.getenv('DB_PASSWORD')
        self.database = os.getenv('DB_NAME')
        
    def validate(self):
        """Validate required configuration"""
        if not all([self.user, self.password, self.database]):
            raise ValueError("Missing required database configuration")


class TwilioConfig:
    """Twilio configuration from environment variables"""
    def __init__(self):
        self.account_sid = os.getenv('TWILIO_ACCOUNT_SID')
        self.auth_token = os.getenv('TWILIO_AUTH_TOKEN')
        self.from_number = os.getenv('TWILIO_FROM_NUMBER')
        
    def validate(self):
        """Validate required configuration"""
        if not all([self.account_sid, self.auth_token, self.from_number]):
            raise ValueError("Missing required Twilio configuration")


class HabitNotifier:
    """Main application class for habit notifications"""
    
    def __init__(self, db_config: DatabaseConfig, twilio_config: TwilioConfig):
        self.db_config = db_config
        self.twilio_config = twilio_config
        self.twilio_client = None
        self.db_connection = None
        
    def connect_db(self):
        """Establish database connection"""
        try:
            self.db_connection = pymysql.connect(
                host=self.db_config.host,
                port=self.db_config.port,
                user=self.db_config.user,
                password=self.db_config.password,
                database=self.db_config.database,
                charset='utf8mb4',
                cursorclass=pymysql.cursors.DictCursor
            )
            logger.info("Database connection established")
        except Exception as e:
            logger.error(f"Failed to connect to database: {e}")
            raise
            
    def connect_twilio(self):
        """Initialize Twilio client"""
        try:
            self.twilio_client = Client(
                self.twilio_config.account_sid,
                self.twilio_config.auth_token
            )
            logger.info("Twilio client initialized")
        except Exception as e:
            logger.error(f"Failed to initialize Twilio client: {e}")
            raise
            
    def get_users_needing_notification(self) -> List[Dict]:
        """
        Query database for users who need notifications based on their habit logs.
        
        Returns:
            List of user dictionaries with user info and notification reason
        """
        users_to_notify = []
        
        try:
            with self.db_connection.cursor() as cursor:
                # Get all users
                cursor.execute("""
                    SELECT id, name, phone_number, time_zone
                    FROM User
                    WHERE phone_number IS NOT NULL AND phone_number != ''
                """)
                users = cursor.fetchall()
                
                logger.info(f"Found {len(users)} users with phone numbers")
                
                for user in users:
                    user_id = user['id']
                    user_tz = pytz.timezone(user['time_zone']) if user.get('time_zone') else pytz.UTC
                    
                    # Get current time in user's timezone
                    now = datetime.now(user_tz)
                    
                    # Calculate day boundaries in user's timezone
                    today_start = now.replace(hour=0, minute=0, second=0, microsecond=0)
                    yesterday_start = today_start - timedelta(days=1)
                    two_days_ago_start = today_start - timedelta(days=2)
                    
                    # Convert to UTC for database query
                    today_start_utc = today_start.astimezone(pytz.UTC)
                    yesterday_start_utc = yesterday_start.astimezone(pytz.UTC)
                    two_days_ago_start_utc = two_days_ago_start.astimezone(pytz.UTC)
                    
                    # Query logs for this user over the past 2 days
                    cursor.execute("""
                        SELECT 
                            DATE(CONVERT_TZ(l.created_at, '+00:00', %s)) as log_date,
                            COUNT(*) as log_count
                        FROM Log l
                        INNER JOIN Habit h ON l.habit_id = h.id
                        WHERE h.user_id = %s
                        AND l.created_at >= %s
                        GROUP BY DATE(CONVERT_TZ(l.created_at, '+00:00', %s))
                        ORDER BY log_date DESC
                    """, (
                        user_tz.zone,
                        user_id,
                        two_days_ago_start_utc,
                        user_tz.zone
                    ))
                    
                    logs = cursor.fetchall()
                    
                    # Create a set of dates with logs
                    logged_dates = {log['log_date'] for log in logs}
                    
                    yesterday_date = yesterday_start.date()
                    two_days_ago_date = two_days_ago_start.date()
                    
                    has_yesterday_log = yesterday_date in logged_dates
                    has_two_days_ago_log = two_days_ago_date in logged_dates
                    
                    # Determine if notification is needed
                    should_notify = False
                    reason = ""
                    
                    if not has_yesterday_log and not has_two_days_ago_log:
                        should_notify = True
                        reason = "No logs in the past 2 days"
                    elif has_two_days_ago_log and not has_yesterday_log:
                        should_notify = True
                        reason = "Logged 2 days ago but not yesterday"
                    
                    if should_notify:
                        users_to_notify.append({
                            'user_id': user_id,
                            'name': user['name'],
                            'phone_number': user['phone_number'],
                            'reason': reason
                        })
                        logger.info(f"User {user['name']} (ID: {user_id}) needs notification: {reason}")
                
                logger.info(f"Total users needing notification: {len(users_to_notify)}")
                return users_to_notify
                
        except Exception as e:
            logger.error(f"Error querying database: {e}")
            raise
            
    def send_notification(self, user: Dict) -> bool:
        """
        Send SMS notification to a user via Twilio.
        
        Args:
            user: Dictionary containing user information
            
        Returns:
            True if notification sent successfully, False otherwise
        """
        try:
            message_body = (
                f"Hi {user['name']}! 👋\n\n"
                f"We noticed you haven't logged your habits recently. "
                f"Don't break your streak! Log your progress today to stay on track. 💪"
            )
            
            message = self.twilio_client.messages.create(
                body=message_body,
                from_=self.twilio_config.from_number,
                to=user['phone_number']
            )
            
            logger.info(
                f"Notification sent to {user['name']} ({user['phone_number']}). "
                f"Message SID: {message.sid}"
            )
            return True
            
        except TwilioRestException as e:
            logger.error(
                f"Twilio error sending notification to {user['name']}: "
                f"{e.code} - {e.msg}"
            )
            return False
        except Exception as e:
            logger.error(f"Error sending notification to {user['name']}: {e}")
            return False
            
    def run(self):
        """Main execution method"""
        success_count = 0
        failure_count = 0
        
        try:
            logger.info("Starting habit notification job")
            
            # Connect to services
            self.connect_db()
            self.connect_twilio()
            
            # Get users needing notifications
            users = self.get_users_needing_notification()
            
            if not users:
                logger.info("No users need notifications at this time")
                return
            
            # Send notifications
            for user in users:
                if self.send_notification(user):
                    success_count += 1
                else:
                    failure_count += 1
                    
            logger.info(
                f"Job completed. Notifications sent: {success_count}, "
                f"Failed: {failure_count}"
            )
            
        except Exception as e:
            logger.error(f"Job failed with error: {e}")
            sys.exit(1)
        finally:
            if self.db_connection:
                self.db_connection.close()
                logger.info("Database connection closed")
                
        # Exit with error code if any notifications failed
        if failure_count > 0:
            sys.exit(1)


def main():
    """Entry point for the application"""
    try:
        # Load and validate configuration
        db_config = DatabaseConfig()
        db_config.validate()
        
        twilio_config = TwilioConfig()
        twilio_config.validate()
        
        # Run the notifier
        notifier = HabitNotifier(db_config, twilio_config)
        notifier.run()
        
    except ValueError as e:
        logger.error(f"Configuration error: {e}")
        sys.exit(1)
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
