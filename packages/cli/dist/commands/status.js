"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.status = status;
const chalk_1 = __importDefault(require("chalk"));
const ora_1 = __importDefault(require("ora"));
const net_1 = require("net");
async function status(options) {
    console.log(chalk_1.default.cyan('\n🔍 Checking Chaos-Proxy status...\n'));
    // Parse Redis URL
    const redisUrl = new URL(options.redis);
    const host = redisUrl.hostname;
    const port = parseInt(redisUrl.port) || 6379;
    // Check Redis
    const redisSpinner = (0, ora_1.default)('Checking Redis connection...').start();
    try {
        await checkConnection(host, port);
        redisSpinner.succeed(chalk_1.default.green(`Redis: Connected (${host}:${port})`));
    }
    catch {
        redisSpinner.fail(chalk_1.default.red(`Redis: Not reachable (${host}:${port})`));
    }
    // Check Sentinel (default port)
    const sentinelSpinner = (0, ora_1.default)('Checking Sentinel proxy...').start();
    try {
        await checkConnection('localhost', 8080);
        sentinelSpinner.succeed(chalk_1.default.green('Sentinel: Running (localhost:8080)'));
    }
    catch {
        sentinelSpinner.fail(chalk_1.default.yellow('Sentinel: Not running (localhost:8080)'));
    }
    // Check Dashboard
    const dashboardSpinner = (0, ora_1.default)('Checking Dashboard...').start();
    try {
        await checkConnection('localhost', 3000);
        dashboardSpinner.succeed(chalk_1.default.green('Dashboard: Running (localhost:3000)'));
    }
    catch {
        dashboardSpinner.fail(chalk_1.default.yellow('Dashboard: Not running (localhost:3000)'));
    }
    console.log('');
}
function checkConnection(host, port) {
    return new Promise((resolve, reject) => {
        const socket = (0, net_1.createConnection)({ host, port }, () => {
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
//# sourceMappingURL=status.js.map