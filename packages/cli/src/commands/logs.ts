import chalk from 'chalk';
import { spawn } from 'child_process';
import { existsSync } from 'fs';

interface LogsOptions {
    service?: string;
    follow: boolean;
}

export async function logs(options: LogsOptions): Promise<void> {
    console.log(chalk.cyan('\n📜 Fetching Chaos-Proxy logs...\n'));

    // Check if docker-compose.yml exists
    if (!existsSync('docker-compose.yml')) {
        console.log(chalk.red('❌ docker-compose.yml not found!'));
        console.log(chalk.yellow('   This command only works with Docker setup.'));
        process.exit(1);
    }

    const args = ['logs'];
    if (options.follow) {
        args.push('-f');
    }

    const validServices = ['sentinel', 'brain', 'dashboard', 'redis'];

    if (options.service) {
        if (!validServices.includes(options.service)) {
            console.log(chalk.red(`❌ Invalid service: ${options.service}`));
            console.log(chalk.yellow(`   Available services: ${validServices.join(', ')}`));
            process.exit(1);
        }
        args.push(options.service);
    }

    // Use shell: true because docker-compose might not be in PATH directly on some systems
    // But since we validated inputs, it's safer.
    const child = spawn('docker-compose', args, {
        stdio: 'inherit',
        shell: true
    });

    child.on('error', (err) => {
        console.error(chalk.red('Failed to run docker-compose logs'));
        console.error(err);
    });
}
