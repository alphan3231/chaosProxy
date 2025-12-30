import { NextRequest, NextResponse } from 'next/server';
import redis from '@/lib/redis';

export async function GET() {
    try {
        const settings = await redis.hgetall('chaos:settings');

        // Convert strings to proper types
        const parsedSettings = {
            latency_enabled: settings.latency_enabled === 'true',
            latency_min: parseInt(settings.latency_min || '500'),
            latency_max: parseInt(settings.latency_max || '1500'),
            failure_enabled: settings.failure_enabled === 'true',
            failure_rate: parseInt(settings.failure_rate || '10'),
        };

        return NextResponse.json(parsedSettings);
    } catch (error) {
        console.error('Failed to fetch chaos settings:', error);
        return NextResponse.json({ error: 'Failed to fetch settings' }, { status: 500 });
    }
}

export async function POST(request: NextRequest) {
    try {
        const body = await request.json();

        // Validate and prepare data for Redis
        const settings = {
            latency_enabled: String(body.latency_enabled),
            latency_min: String(body.latency_min),
            latency_max: String(body.latency_max),
            failure_enabled: String(body.failure_enabled),
            failure_rate: String(body.failure_rate),
        };

        await redis.hset('chaos:settings', settings);

        return NextResponse.json({ success: true, settings: body });
    } catch (error) {
        console.error('Failed to update chaos settings:', error);
        return NextResponse.json({ error: 'Failed to update settings' }, { status: 500 });
    }
}
