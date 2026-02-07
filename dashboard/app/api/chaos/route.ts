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

        // Validation
        const latencyMin = parseInt(body.latency_min);
        const latencyMax = parseInt(body.latency_max);
        const failureRate = parseInt(body.failure_rate);

        if (isNaN(latencyMin) || latencyMin < 0) {
            return NextResponse.json({ error: 'Invalid latency_min' }, { status: 400 });
        }
        if (isNaN(latencyMax) || latencyMax < 0) {
            return NextResponse.json({ error: 'Invalid latency_max' }, { status: 400 });
        }
        if (latencyMax < latencyMin) {
             return NextResponse.json({ error: 'latency_max must be greater than or equal to latency_min' }, { status: 400 });
        }
        if (isNaN(failureRate) || failureRate < 0 || failureRate > 100) {
            return NextResponse.json({ error: 'failure_rate must be between 0 and 100' }, { status: 400 });
        }

        // Validate and prepare data for Redis
        const settings = {
            latency_enabled: body.latency_enabled === true || body.latency_enabled === 'true' ? 'true' : 'false',
            latency_min: String(latencyMin),
            latency_max: String(latencyMax),
            failure_enabled: body.failure_enabled === true || body.failure_enabled === 'true' ? 'true' : 'false',
            failure_rate: String(failureRate),
        };

        await redis.hset('chaos:settings', settings);

        return NextResponse.json({ success: true, settings: body });
    } catch (error) {
        console.error('Failed to update chaos settings:', error);
        return NextResponse.json({ error: 'Failed to update settings' }, { status: 500 });
    }
}
