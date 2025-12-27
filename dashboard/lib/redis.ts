import Redis from 'ioredis';

const redisUrl = process.env.REDIS_URL || 'redis://localhost:6379';
const redisPassword = process.env.REDIS_PASSWORD || undefined;

const redis = new Redis(redisUrl, {
    password: redisPassword,
    maxRetriesPerRequest: 3,
});

export default redis;
