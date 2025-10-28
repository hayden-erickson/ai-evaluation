import logging
from twilio.rest import Client

logger = logging.getLogger(__name__)

_client: Client | None = None

def init_client(account_sid: str, auth_token: str) -> None:
    global _client
    if _client is not None:
        return
    _client = Client(account_sid, auth_token)


def send_sms(from_number: str, to_number: str, body: str) -> None:
    if _client is None:
        raise RuntimeError("Twilio client not initialized")
    try:
        _client.messages.create(from_=from_number, to=to_number, body=body)
        logger.info("sms_sent", extra={"to": to_number})
    except Exception as e:
        logger.error("sms_error", extra={"to": to_number, "error": str(e)})
        raise
