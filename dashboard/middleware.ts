import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function middleware(request: NextRequest) {
    const path = request.nextUrl.pathname;

    // Define public paths (static assets, images, etc.)
    // We want to protect root "/" and "/api/*"
    const isPublic = path.startsWith('/_next') ||
                     path.startsWith('/static') ||
                     path.includes('.'); // Simple check for files like favicon.ico, logo.svg

    if (isPublic) {
        return NextResponse.next();
    }

    const authHeader = request.headers.get('authorization');

    if (!authHeader) {
        return new NextResponse(JSON.stringify({ error: 'Authentication required' }), {
            status: 401,
            headers: {
                'WWW-Authenticate': 'Basic realm="Chaos Dashboard"',
                'Content-Type': 'application/json',
            },
        });
    }

    const [scheme, encoded] = authHeader.split(' ');

    if (scheme !== 'Basic' || !encoded) {
        return new NextResponse(JSON.stringify({ error: 'Invalid authentication scheme' }), {
            status: 401,
            headers: { 'Content-Type': 'application/json' },
        });
    }

    const decoded = Buffer.from(encoded, 'base64').toString('utf-8');
    const [username, password] = decoded.split(':');

    const validUser = process.env.DASHBOARD_USER || 'admin';
    const validPass = process.env.DASHBOARD_PASSWORD || 'chaos123';

    if (username !== validUser || password !== validPass) {
        return new NextResponse(JSON.stringify({ error: 'Invalid credentials' }), {
            status: 403,
            headers: { 'Content-Type': 'application/json' },
        });
    }

    return NextResponse.next();
}

export const config = {
    matcher: ['/((?!_next/static|_next/image|favicon.ico).*)'],
};
