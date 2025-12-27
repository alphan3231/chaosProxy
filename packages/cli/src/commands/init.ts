import chalk from 'chalk';
import ora from 'ora';
import inquirer from 'inquirer';
import { writeFileSync, mkdirSync, existsSync } from 'fs';
import { join } from 'path';

interface InitOptions {
    docker: boolean;
}

export async function init(options: InitOptions): Promise<void> {
    console.log(chalk.cyan('\n🚀 Initializing Chaos-Proxy project...\n'));

    const answers = await inquirer.prompt([
        {
            type: 'input',
            name: 'projectName',
            message: 'Project name:',
            default: 'my-chaos-proxy',
        },
        {
            type: 'input',
            name: 'targetUrl',
            message: 'Target backend URL:',
            default: 'http://localhost:3000',
        },
        {
            type: 'input',
            name: 'proxyPort',
            message: 'Proxy port:',
            default: '8080',
        },
        {
            type: 'confirm',
            name: 'useDocker',
            message: 'Use Docker for deployment?',
            default: options.docker,
        },
        {
            type: 'password',
            name: 'redisPassword',
            message: 'Redis password (leave empty for no password):',
            default: '',
        },
    ]);

    const spinner = ora('Creating project files...').start();

    try {
        const projectDir = join(process.cwd(), answers.projectName);

        if (!existsSync(projectDir)) {
            mkdirSync(projectDir, { recursive: true });
        }

        // Create .env file
        const envContent = `# Chaos-Proxy Configuration
PORT=${answers.proxyPort}
TARGET_URL=${answers.targetUrl}
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=${answers.redisPassword}
APP_ENV=development

# Dashboard
DASHBOARD_USER=admin
DASHBOARD_PASSWORD=chaos123
`;

        writeFileSync(join(projectDir, '.env'), envContent);

        // Create docker-compose.yml if Docker is selected
        if (answers.useDocker) {
            const dockerCompose = `version: '3.8'

services:
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    command: redis-server ${answers.redisPassword ? `--requirepass ${answers.redisPassword}` : ''}

  sentinel:
    image: elliot/chaos-proxy:sentinel
    ports:
      - "${answers.proxyPort}:8080"
    environment:
      - TARGET_URL=${answers.targetUrl}
      - REDIS_ADDR=redis:6379
      - REDIS_PASSWORD=${answers.redisPassword}
    depends_on:
      - redis

  brain:
    image: elliot/chaos-proxy:brain
    environment:
      - REDIS_ADDR=redis
      - REDIS_PORT=6379
      - REDIS_PASSWORD=${answers.redisPassword}
    depends_on:
      - redis

  dashboard:
    image: elliot/chaos-proxy:dashboard
    ports:
      - "3000:3000"
    environment:
      - REDIS_URL=redis://redis:6379
      - REDIS_PASSWORD=${answers.redisPassword}
    depends_on:
      - redis
`;

            writeFileSync(join(projectDir, 'docker-compose.yml'), dockerCompose);
        }

        spinner.succeed(chalk.green('Project created successfully!'));

        console.log(chalk.cyan('\n📁 Project structure:'));
        console.log(`   ${answers.projectName}/`);
        console.log('   ├── .env');
        if (answers.useDocker) {
            console.log('   └── docker-compose.yml');
        }

        console.log(chalk.yellow('\n📋 Next steps:'));
        console.log(`   cd ${answers.projectName}`);
        if (answers.useDocker) {
            console.log('   docker-compose up -d');
        } else {
            console.log('   # Install Chaos-Proxy from GitHub');
            console.log('   git clone https://github.com/elliot/chaosProxy.git');
        }

        console.log(chalk.green('\n✨ Happy proxying! 👻\n'));

    } catch (error) {
        spinner.fail(chalk.red('Failed to create project'));
        console.error(error);
        process.exit(1);
    }
}
