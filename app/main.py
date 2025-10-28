import logging
import sys

from .config import Config
from . import db, notifier
from .health import start_health_server, stop_health_server
from .logic import run_job


def setup_logging(level: str) -> None:
    logging.basicConfig(stream=sys.stdout, level=getattr(logging, level.upper(), logging.INFO), format="%(asctime)s %(levelname)s %(name)s %(message)s")


def main() -> int:
    cfg = Config()
    cfg.validate()
    setup_logging(cfg.log_level)
    server = start_health_server()
    db.init_pool(
        cfg.mysql_host,
        cfg.mysql_port,
        cfg.mysql_db,
        cfg.mysql_user,
        cfg.mysql_password,
        ssl_disabled=cfg.mysql_ssl_disabled,
        ssl_ca_path=cfg.mysql_ssl_ca_path or None,
    )
    notifier.init_client(cfg.twilio_account_sid, cfg.twilio_auth_token)
    notified, skipped, errors = run_job(cfg.twilio_from)
    logging.getLogger(__name__).info("run_complete", extra={"notified": notified, "skipped": skipped, "errors": errors})
    stop_health_server(server)
    return 0 if errors == 0 else 1


if __name__ == "__main__":
    sys.exit(main())
