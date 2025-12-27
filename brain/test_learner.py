"""
Tests for the Brain Learner module
"""
import unittest
import json
import sys
from unittest.mock import MagicMock, patch

# Mock redis module before importing learner
mock_redis_module = MagicMock()
sys.modules['redis'] = mock_redis_module


class TestLearner(unittest.TestCase):
    """Test cases for Learner class"""

    def setUp(self):
        """Set up test fixtures"""
        self.mock_redis = MagicMock()
        
        # Import after setting up mocks
        from learner import Learner
        self.learner = Learner(self.mock_redis)

    def test_learn_successful_response(self):
        """Test that 200 responses are learned"""
        traffic_data = {
            'method': 'GET',
            'path': '/api/users',
            'status': 200,
            'response_body': '{"users": []}',
            'timestamp': '2025-01-01T00:00:00Z'
        }

        self.learner.learn(traffic_data)

        # Verify Redis set was called
        self.mock_redis.set.assert_called_once()
        call_args = self.mock_redis.set.call_args

        # Check key format
        key = call_args[0][0]
        self.assertEqual(key, 'chaos:ghost:GET:/api/users')

        # Check that request_count was incremented
        self.mock_redis.incr.assert_called_with('chaos:stats:request_count')

    def test_learn_error_response_not_cached(self):
        """Test that 4xx/5xx responses are not cached"""
        traffic_data = {
            'method': 'GET',
            'path': '/api/users',
            'status': 404,
            'response_body': '{"error": "Not found"}',
            'timestamp': '2025-01-01T00:00:00Z'
        }

        self.learner.learn(traffic_data)

        # set should NOT be called for error responses
        self.mock_redis.set.assert_not_called()

        # But counter should still be incremented
        self.mock_redis.incr.assert_called_with('chaos:stats:request_count')

    def test_learn_201_response(self):
        """Test that 201 Created responses are learned"""
        traffic_data = {
            'method': 'POST',
            'path': '/api/users',
            'status': 201,
            'response_body': '{"id": 1, "name": "John"}',
            'timestamp': '2025-01-01T00:00:00Z'
        }

        self.learner.learn(traffic_data)

        # set should be called for 201
        self.mock_redis.set.assert_called_once()

    def test_learn_no_status(self):
        """Test handling of traffic data without status"""
        traffic_data = {
            'method': 'GET',
            'path': '/api/users',
            'response_body': '{"users": []}',
            'timestamp': '2025-01-01T00:00:00Z'
        }

        # Should not raise exception
        self.learner.learn(traffic_data)

        # set should NOT be called without valid status
        self.mock_redis.set.assert_not_called()

    def test_ghost_response_format(self):
        """Test that ghost response has correct format"""
        traffic_data = {
            'method': 'GET',
            'path': '/api/test',
            'status': 200,
            'response_body': '{"test": true}',
            'timestamp': '2025-01-01T00:00:00Z'
        }

        self.learner.learn(traffic_data)

        # Get the value that was set
        call_args = self.mock_redis.set.call_args
        value = json.loads(call_args[0][1])

        self.assertEqual(value['status'], 200)
        self.assertEqual(value['body'], '{"test": true}')
        self.assertIn('timestamp', value)

    def test_ttl_is_set(self):
        """Test that TTL is set on cached responses"""
        traffic_data = {
            'method': 'GET',
            'path': '/api/test',
            'status': 200,
            'response_body': '{}',
            'timestamp': '2025-01-01T00:00:00Z'
        }

        self.learner.learn(traffic_data)

        # Check that ex (expiration) parameter is passed
        call_kwargs = self.mock_redis.set.call_args[1]
        self.assertIn('ex', call_kwargs)
        self.assertEqual(call_kwargs['ex'], 3600)  # 1 hour


class TestLearnerEdgeCases(unittest.TestCase):
    """Edge case tests for Learner"""

    def setUp(self):
        self.mock_redis = MagicMock()
        from learner import Learner
        self.learner = Learner(self.mock_redis)

    def test_empty_response_body(self):
        """Test handling of empty response body"""
        traffic_data = {
            'method': 'DELETE',
            'path': '/api/users/1',
            'status': 204,
            'response_body': '',
            'timestamp': '2025-01-01T00:00:00Z'
        }

        self.learner.learn(traffic_data)
        self.mock_redis.set.assert_called_once()

    def test_special_characters_in_path(self):
        """Test handling of special characters in path"""
        traffic_data = {
            'method': 'GET',
            'path': '/api/users?id=1&name=test',
            'status': 200,
            'response_body': '{}',
            'timestamp': '2025-01-01T00:00:00Z'
        }

        self.learner.learn(traffic_data)

        call_args = self.mock_redis.set.call_args
        key = call_args[0][0]
        self.assertEqual(key, 'chaos:ghost:GET:/api/users?id=1&name=test')


if __name__ == '__main__':
    unittest.main()
