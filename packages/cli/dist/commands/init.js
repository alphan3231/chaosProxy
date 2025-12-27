"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.init = init;
const chalk_1 = __importDefault(require("chalk"));
const ora_1 = __importDefault(require("ora"));
const inquirer_1 = __importDefault(require("inquirer"));
const fs_1 = require("fs");
const path_1 = require("path");
async function init(options) {
    console.log(chalk_1.default.cyan('\n🚀 Initializing Chaos-Proxy project...\n'));
    const answers = await inquirer_1.default.prompt([
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
    const spinner = (0, ora_1.default)('Creating project files...').start();
    try {
        const projectDir = (0, path_1.join)(process.cwd(), answers.projectName);
        if (!(0, fs_1.existsSync)(projectDir)) {
            (0, fs_1.mkdirSync)(projectDir, { recursive: true });
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
        (0, fs_1.writeFileSync)((0, path_1.join)(projectDir, '.env'), envContent);
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
            (0, fs_1.writeFileSync)((0, path_1.join)(projectDir, 'docker-compose.yml'), dockerCompose);
        }
        spinner.succeed(chalk_1.default.green('Project created successfully!'));
        console.log(chalk_1.default.cyan('\n📁 Project structure:'));
        console.log(`   ${answers.projectName}/`);
        console.log('   ├── .env');
        if (answers.useDocker) {
            console.log('   └── docker-compose.yml');
        }
        console.log(chalk_1.default.yellow('\n📋 Next steps:'));
        console.log(`   cd ${answers.projectName}`);
        if (answers.useDocker) {
            console.log('   docker-compose up -d');
        }
        else {
            console.log('   # Install Chaos-Proxy from GitHub');
            console.log('   git clone https://github.com/elliot/chaosProxy.git');
        }
        console.log(chalk_1.default.green('\n✨ Happy proxying! 👻\n'));
    }
    catch (error) {
        spinner.fail(chalk_1.default.red('Failed to create project'));
        console.error(error);
        process.exit(1);
    }
}
//# sourceMappingURL=init.js.map