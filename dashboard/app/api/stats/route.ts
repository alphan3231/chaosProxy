import { NextResponse } from 'next/server';
import redis from '@/lib/redis';

export const dynamic = 'force-dynamic';

export async function GET() {
    try {
        const [totalRequests, ghostCount] = await Promise.all([
            redis.get('chaos:stats:request_count'),
            redis.get('chaos:stats:ghost_count'),
        ]);

        return NextResponse.json({
            totalRequests: parseInt(totalRequests || '0'),
            ghostCount: parseInt(ghostCount || '0'),
            status: 'active', // Should be checked properly in real scenario
        });
    } catch (error) {
        return NextResponse.json({ error: 'Failed to fetch stats' }, { status: 500 });
    }
}
