import { NextResponse } from 'next/server';
import redis from '@/lib/redis';

export const dynamic = 'force-dynamic';

export async function GET() {
    try {
        const [totalRequests, ghostCount] = await Promise.all([
            redis.get('chaos:stats:request_count'),
            redis.get('chaos:stats:ghost_count'),
        ]);

        // Fetch recent logs
        const recentLogsRaw = await redis.lrange('chaos:logs:recent', 0, 19); // Get last 20
        const recentLogs = recentLogsRaw.map(log => JSON.parse(log));

        return NextResponse.json({
            totalRequests: parseInt(totalRequests || '0'),
            ghostCount: parseInt(ghostCount || '0'),
            status: 'active',
            recentLogs
        });
    } catch (error) {
        return NextResponse.json({ error: 'Failed to fetch stats' }, { status: 500 });
    }
}
