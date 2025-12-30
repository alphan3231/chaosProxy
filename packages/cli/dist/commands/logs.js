"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.logs = logs;
const chalk_1 = __importDefault(require("chalk"));
const child_process_1 = require("child_process");
const fs_1 = require("fs");
async function logs(options) {
    console.log(chalk_1.default.cyan('\n📜 Fetching Chaos-Proxy logs...\n'));
    // Check if docker-compose.yml exists
    if (!(0, fs_1.existsSync)('docker-compose.yml')) {
        console.log(chalk_1.default.red('❌ docker-compose.yml not found!'));
        console.log(chalk_1.default.yellow('   This command only works with Docker setup.'));
        process.exit(1);
    }
    const args = ['logs'];
    if (options.follow) {
        args.push('-f');
    }
    if (options.service) {
        args.push(options.service);
    }
    const child = (0, child_process_1.spawn)('docker-compose', args, {
        stdio: 'inherit',
        shell: true
    });
    child.on('error', (err) => {
        console.error(chalk_1.default.red('Failed to run docker-compose logs'));
        console.error(err);
    });
}
//# sourceMappingURL=logs.js.map