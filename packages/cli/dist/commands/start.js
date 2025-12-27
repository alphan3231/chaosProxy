"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.start = start;
const chalk_1 = __importDefault(require("chalk"));
const ora_1 = __importDefault(require("ora"));
const child_process_1 = require("child_process");
const fs_1 = require("fs");
async function start(options) {
    console.log(chalk_1.default.cyan('\n🚀 Starting Chaos-Proxy...\n'));
    if (options.docker) {
        // Check if docker-compose.yml exists
        if (!(0, fs_1.existsSync)('docker-compose.yml')) {
            console.log(chalk_1.default.red('❌ docker-compose.yml not found!'));
            console.log(chalk_1.default.yellow('   Run "chaos-proxy init --docker" first.'));
            process.exit(1);
        }
        const spinner = (0, ora_1.default)('Starting Docker containers...').start();
        try {
            (0, child_process_1.execSync)('docker-compose up -d', { stdio: 'pipe' });
            spinner.succeed(chalk_1.default.green('Docker containers started!'));
            console.log(chalk_1.default.cyan('\n📊 Services:'));
            console.log('   • Proxy:     http://localhost:8080');
            console.log('   • Dashboard: http://localhost:3000');
            console.log('   • Redis:     localhost:6379');
            console.log(chalk_1.default.yellow('\n💡 Tips:'));
            console.log('   • View logs:  docker-compose logs -f');
            console.log('   • Stop:       docker-compose down');
            console.log('   • Status:     chaos-proxy status\n');
        }
        catch (error) {
            spinner.fail(chalk_1.default.red('Failed to start containers'));
            console.error(error);
            process.exit(1);
        }
    }
    else {
        console.log(chalk_1.default.yellow('⚠️  Non-Docker mode requires manual setup.\n'));
        console.log('Please run the following in separate terminals:');
        console.log(chalk_1.default.cyan('\n1. Start Redis:'));
        console.log('   docker run -d -p 6379:6379 redis:7-alpine');
        console.log(chalk_1.default.cyan('\n2. Start Sentinel (from chaosProxy repo):'));
        console.log('   go run cmd/sentinel/main.go');
        console.log(chalk_1.default.cyan('\n3. Start Brain (from chaosProxy repo):'));
        console.log('   cd brain && python main.py');
        console.log(chalk_1.default.cyan('\n4. Start Dashboard (from chaosProxy repo):'));
        console.log('   cd dashboard && npm run dev');
        console.log(chalk_1.default.yellow('\nOr use: chaos-proxy start --docker\n'));
    }
}
//# sourceMappingURL=start.js.map