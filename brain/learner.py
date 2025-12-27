import json
import logging
from redis import Redis

logger = logging.getLogger(__name__)

class Learner:
    def __init__(self, redis_client: Redis):
        self.redis = redis_client
        self.ttl = 3600 # 1 hour for now

    def learn(self, traffic_data: dict):
        """
        Analyzes traffic and updates the Ghost Model.
        For MVP: Simply caches the last successful response for a given Method+Path.
        """
        method = traffic_data.get('method')
        path = traffic_data.get('path')
        status = traffic_data.get('status')
        
        # Only learn from successful responses for now
        if status and 200 <= status < 300:
            self._save_ghost_response(method, path, traffic_data)

    def _save_ghost_response(self, method, path, data):
        key = f"chaos:ghost:{method}:{path}"
        
        # We store the structure that the Ghost (Go) will eventually read
        ghost_response = {
            "status": data.get('status'),
            "body": data.get('response_body'),
            "headers": {}, # We didn't capture headers in MVP yet
            "timestamp": data.get('timestamp')
        }
        
        try:
            self.redis.set(key, json.dumps(ghost_response), ex=self.ttl)
            logger.info(f"🧠 Learned pattern for: {method} {path}")
        except Exception as e:
            logger.error(f"Failed to save ghost response: {e}")
