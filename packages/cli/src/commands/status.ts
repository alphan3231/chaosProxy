import chalk from 'chalk';
import ora from 'ora';
import { createConnection } from 'net';

interface StatusOptions {
    redis: string;
}

export async function status(options: StatusOptions): Promise<void> {
    console.log(chalk.cyan('\n🔍 Checking Chaos-Proxy status...\n'));

    // Parse Redis URL
    const redisUrl = new URL(options.redis);
    const host = redisUrl.hostname;
    const port = parseInt(redisUrl.port) || 6379;

    // Check Redis
    const redisSpinner = ora('Checking Redis connection...').start();

    try {
        await checkConnection(host, port);
        redisSpinner.succeed(chalk.green(`Redis: Connected (${host}:${port})`));
    } catch {
        redisSpinner.fail(chalk.red(`Redis: Not reachable (${host}:${port})`));
    }

    // Check Sentinel (default port)
    const sentinelSpinner = ora('Checking Sentinel proxy...').start();

    try {
        await checkConnection('localhost', 8080);
        sentinelSpinner.succeed(chalk.green('Sentinel: Running (localhost:8080)'));
    } catch {
        sentinelSpinner.fail(chalk.yellow('Sentinel: Not running (localhost:8080)'));
    }

    // Check Dashboard
    const dashboardSpinner = ora('Checking Dashboard...').start();

    try {
        await checkConnection('localhost', 3000);
        dashboardSpinner.succeed(chalk.green('Dashboard: Running (localhost:3000)'));
    } catch {
        dashboardSpinner.fail(chalk.yellow('Dashboard: Not running (localhost:3000)'));
    }

    console.log('');
}

function checkConnection(host: string, port: number): Promise<void> {
    return new Promise((resolve, reject) => {
        const socket = createConnection({ host, port }, () => {
            socket.destroy();
            resolve();
        });

        socket.on('error', () => {
            socket.destroy();
            reject(new Error('Connection failed'));
        });

        socket.setTimeout(2000, () => {
            socket.destroy();
            reject(new Error('Connection timeout'));
        });
    });
}
