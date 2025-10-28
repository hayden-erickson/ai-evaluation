import os

class Config:
    def __init__(self) -> None:
        self.log_level = os.getenv("LOG_LEVEL", "INFO")
        self.mysql_host = os.getenv("MYSQL_HOST", "")
        self.mysql_port = int(os.getenv("MYSQL_PORT", "3306"))
        self.mysql_db = os.getenv("MYSQL_DB", "")
        self.mysql_user = os.getenv("MYSQL_USER", "")
        self.mysql_password = os.getenv("MYSQL_PASSWORD", "")
        self.mysql_ssl_disabled = os.getenv("MYSQL_SSL_DISABLED", "false").lower() == "true"
        self.mysql_ssl_ca_path = os.getenv("MYSQL_SSL_CA_PATH", "")
        self.twilio_account_sid = os.getenv("TWILIO_ACCOUNT_SID", "")
        self.twilio_auth_token = os.getenv("TWILIO_AUTH_TOKEN", "")
        self.twilio_from = os.getenv("TWILIO_FROM", "")

    def validate(self) -> None:
        missing = []
        if not self.mysql_host:
            missing.append("MYSQL_HOST")
        if not self.mysql_db:
            missing.append("MYSQL_DB")
        if not self.mysql_user:
            missing.append("MYSQL_USER")
        if not self.mysql_password:
            missing.append("MYSQL_PASSWORD")
        if not self.twilio_account_sid:
            missing.append("TWILIO_ACCOUNT_SID")
        if not self.twilio_auth_token:
            missing.append("TWILIO_AUTH_TOKEN")
        if not self.twilio_from:
            missing.append("TWILIO_FROM")
        if missing:
            raise ValueError("Missing required configuration: " + ", ".join(missing))
