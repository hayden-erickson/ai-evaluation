import logging
from datetime import datetime, timedelta, timezone
from zoneinfo import ZoneInfo
from typing import Tuple

from . import db, notifier

logger = logging.getLogger(__name__)


def _compute_window(tz_name: str) -> Tuple[datetime, datetime, datetime, datetime]:
    now_utc = datetime.now(timezone.utc)
    try:
        zone = ZoneInfo(tz_name)
    except Exception:
        logger.warning("invalid_timezone", extra={"tz": tz_name})
        zone = ZoneInfo("UTC")
    local_now = now_utc.astimezone(zone)
    local_today = local_now.date()
    local_yesterday = local_today - timedelta(days=1)
    local_day_before = local_today - timedelta(days=2)
    start_local_day_before = datetime.combine(local_day_before, datetime.min.time(), tzinfo=zone)
    start_local_today = datetime.combine(local_today, datetime.min.time(), tzinfo=zone)
    start_utc_day_before = start_local_day_before.astimezone(timezone.utc)
    start_utc_today = start_local_today.astimezone(timezone.utc)
    start_local_yesterday = datetime.combine(local_yesterday, datetime.min.time(), tzinfo=zone)
    start_utc_yesterday = start_local_yesterday.astimezone(timezone.utc)
    return start_utc_day_before, start_utc_yesterday, start_utc_today, local_yesterday


def _classify_days(tz_name: str, created_list):
    try:
        z = ZoneInfo(tz_name)
    except Exception:
        logger.warning("invalid_timezone", extra={"tz": tz_name})
        z = ZoneInfo("UTC")
    days = set()
    for dt in created_list:
        if dt.tzinfo is None:
            dt = dt.replace(tzinfo=timezone.utc)
        local_date = dt.astimezone(z).date()
        days.add(local_date)
    return days


def run_job(twilio_from: str) -> Tuple[int, int, int]:
    users = db.fetch_users_with_phones()
    notified = 0
    skipped = 0
    errors = 0
    for user_id, tz_name, phone in users:
        try:
            s1, s2, s3, local_yesterday = _compute_window(tz_name)
            logs = db.fetch_user_logs_between(user_id, s1, s3)
            days = _classify_days(tz_name, logs)
            day_before = (local_yesterday - timedelta(days=1))
            has_day_before = day_before in days
            has_yesterday = local_yesterday in days
            should_notify = False
            reason = ""
            if not has_day_before and not has_yesterday:
                should_notify = True
                reason = "no_logs_two_days"
            elif has_day_before and not has_yesterday:
                should_notify = True
                reason = "missed_yesterday"
            if should_notify:
                body = "Reminder: no recent habit logs. Open the app to log your progress."
                notifier.send_sms(twilio_from, phone, body)
                notified += 1
                logger.info("notified", extra={"user_id": user_id, "reason": reason})
            else:
                skipped += 1
        except Exception as e:
            errors += 1
            logger.error("user_error", extra={"user_id": user_id, "error": str(e)})
    return notified, skipped, errors
