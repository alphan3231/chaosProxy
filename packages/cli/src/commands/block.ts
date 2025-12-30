import chalk from 'chalk';
import ora from 'ora';
import { createClient } from 'redis';

interface BlockOptions {
    redis: string;
}

export async function block(ip: string, options: BlockOptions): Promise<void> {
    const spinner = ora(`Blocking IP: ${ip}...`).start();

    try {
        const client = createClient({ url: options.redis });
        await client.connect();

        await client.sAdd('chaos:settings:blocked_ips', ip);
        await client.quit();

        spinner.succeed(chalk.green(`Blocked IP: ${ip}`));
    } catch (error) {
        spinner.fail(chalk.red(`Failed to block IP: ${ip}`));
        console.error(error);
    }
}

export async function unblock(ip: string, options: BlockOptions): Promise<void> {
    const spinner = ora(`Unblocking IP: ${ip}...`).start();

    try {
        const client = createClient({ url: options.redis });
        await client.connect();

        await client.sRem('chaos:settings:blocked_ips', ip);
        await client.quit();

        spinner.succeed(chalk.green(`Unblocked IP: ${ip}`));
    } catch (error) {
        spinner.fail(chalk.red(`Failed to unblock IP: ${ip}`));
        console.error(error);
    }
}

export async function listBlocked(options: BlockOptions): Promise<void> {
    const spinner = ora('Fetching blocked IPs...').start();

    try {
        const client = createClient({ url: options.redis });
        await client.connect();

        const ips = await client.sMembers('chaos:settings:blocked_ips');
        await client.quit();

        if (ips.length === 0) {
            spinner.info(chalk.yellow('No IPs are currently blocked.'));
            return;
        }

        spinner.stop();
        console.log(chalk.cyan('\n🚫 Blocked IPs:'));
        ips.forEach(ip => console.log(`   - ${ip}`));
        console.log('');

    } catch (error) {
        spinner.fail(chalk.red('Failed to fetch blocked IPs'));
        console.error(error);
    }
}
