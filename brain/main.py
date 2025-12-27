import os
import json
import redis
import logging
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Configuration
REDIS_ADDR = os.getenv("REDIS_ADDR", "localhost")
REDIS_PORT = int(os.getenv("REDIS_PORT", 6379))
CHANNEL = "chaos:traffic"

# Setup Logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - [THE BRAIN] - %(message)s',
    datefmt='%H:%M:%S'
)
logger = logging.getLogger(__name__)

def main():
    logger.info(f"Connecting to Redis at {REDIS_ADDR}:{REDIS_PORT}...")
    
    try:
        r = redis.Redis(host=REDIS_ADDR, port=REDIS_PORT, decode_responses=True)
        pubsub = r.pubsub()
        pubsub.subscribe(CHANNEL)
        
        logger.info(f"Listening on channel: {CHANNEL}")
        logger.info("Ready to learn...")

        for message in pubsub.listen():
            if message['type'] == 'message':
                try:
                    data = json.loads(message['data'])
                    process_traffic(data)
                except json.JSONDecodeError:
                    logger.error("Failed to decode JSON message")
                except Exception as e:
                    logger.error(f"Error processing message: {e}")
                    
    except redis.ConnectionError:
        logger.error("Could not connect to Redis. Is it running?")
    except KeyboardInterrupt:
        logger.info("Shutting down...")

def process_traffic(data):
    """
    Placeholder for the AI logic.
    For now, just acknowledge we saw the traffic.
    """
    method = data.get('method')
    path = data.get('path')
    status = data.get('status')
    logger.info(f"Captured: {method} {path} -> {status}")

if __name__ == "__main__":
    main()
