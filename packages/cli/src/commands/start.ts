import chalk from 'chalk';
import ora from 'ora';
import { spawn, execSync } from 'child_process';
import { existsSync } from 'fs';

interface StartOptions {
    docker: boolean;
}

export async function start(options: StartOptions): Promise<void> {
    console.log(chalk.cyan('\n🚀 Starting Chaos-Proxy...\n'));

    if (options.docker) {
        // Check if docker-compose.yml exists
        if (!existsSync('docker-compose.yml')) {
            console.log(chalk.red('❌ docker-compose.yml not found!'));
            console.log(chalk.yellow('   Run "chaos-proxy init --docker" first.'));
            process.exit(1);
        }

        const spinner = ora('Starting Docker containers...').start();

        try {
            execSync('docker-compose up -d', { stdio: 'pipe' });
            spinner.succeed(chalk.green('Docker containers started!'));

            console.log(chalk.cyan('\n📊 Services:'));
            console.log('   • Proxy:     http://localhost:8080');
            console.log('   • Dashboard: http://localhost:3000');
            console.log('   • Redis:     localhost:6379');

            console.log(chalk.yellow('\n💡 Tips:'));
            console.log('   • View logs:  docker-compose logs -f');
            console.log('   • Stop:       docker-compose down');
            console.log('   • Status:     chaos-proxy status\n');

        } catch (error) {
            spinner.fail(chalk.red('Failed to start containers'));
            console.error(error);
            process.exit(1);
        }
    } else {
        console.log(chalk.yellow('⚠️  Non-Docker mode requires manual setup.\n'));
        console.log('Please run the following in separate terminals:');
        console.log(chalk.cyan('\n1. Start Redis:'));
        console.log('   docker run -d -p 6379:6379 redis:7-alpine');
        console.log(chalk.cyan('\n2. Start Sentinel (from chaosProxy repo):'));
        console.log('   go run cmd/sentinel/main.go');
        console.log(chalk.cyan('\n3. Start Brain (from chaosProxy repo):'));
        console.log('   cd brain && python main.py');
        console.log(chalk.cyan('\n4. Start Dashboard (from chaosProxy repo):'));
        console.log('   cd dashboard && npm run dev');
        console.log(chalk.yellow('\nOr use: chaos-proxy start --docker\n'));
    }
}
