import logging
from typing import List, Tuple
from datetime import datetime

import mysql.connector
from mysql.connector import pooling

logger = logging.getLogger(__name__)

_pool: pooling.MySQLConnectionPool | None = None

def init_pool(host: str, port: int, database: str, user: str, password: str, ssl_disabled: bool = False, ssl_ca_path: str | None = None) -> None:
    global _pool
    if _pool is not None:
        return
    conn_args = {
        "host": host,
        "port": port,
        "database": database,
        "user": user,
        "password": password,
        "autocommit": False,
    }
    if ssl_disabled:
        conn_args["ssl_disabled"] = True
    elif ssl_ca_path:
        try:
            import os
            if os.path.exists(ssl_ca_path):
                conn_args["ssl_ca"] = ssl_ca_path
                conn_args["ssl_verify_cert"] = True
        except Exception:
            pass
    _pool = pooling.MySQLConnectionPool(pool_name="habit_pool", pool_size=8, **conn_args)


def _get_conn():
    if _pool is None:
        raise RuntimeError("DB pool not initialized")
    return _pool.get_connection()


def fetch_users_with_phones() -> List[Tuple[int, str, str]]:
    conn = _get_conn()
    try:
        cur = conn.cursor(dictionary=True)
        cur.execute(
            """
            SELECT `id` AS user_id, `time_zone`, `phone_number`
            FROM `User`
            WHERE `phone_number` IS NOT NULL AND `phone_number` <> ''
            """
        )
        rows = cur.fetchall()
        return [(r["user_id"], r["time_zone"], r["phone_number"]) for r in rows]
    finally:
        conn.close()


def fetch_user_logs_between(user_id: int, start_utc: datetime, end_utc: datetime) -> List[datetime]:
    conn = _get_conn()
    try:
        cur = conn.cursor()
        cur.execute(
            """
            SELECT l.`created_at`
            FROM `Log` l
            JOIN `Habit` h ON h.`id` = l.`habit_id`
            WHERE h.`user_id` = %s AND l.`created_at` >= %s AND l.`created_at` < %s
            """,
            (
                user_id,
                start_utc.replace(tzinfo=None),
                end_utc.replace(tzinfo=None),
            ),
        )
        res = [row[0] for row in cur.fetchall()]
        return res
    finally:
        conn.close()
