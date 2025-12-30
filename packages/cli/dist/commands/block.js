"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.block = block;
exports.unblock = unblock;
exports.listBlocked = listBlocked;
const chalk_1 = __importDefault(require("chalk"));
const ora_1 = __importDefault(require("ora"));
const redis_1 = require("redis");
async function block(ip, options) {
    const spinner = (0, ora_1.default)(`Blocking IP: ${ip}...`).start();
    try {
        const client = (0, redis_1.createClient)({ url: options.redis });
        await client.connect();
        await client.sAdd('chaos:settings:blocked_ips', ip);
        await client.quit();
        spinner.succeed(chalk_1.default.green(`Blocked IP: ${ip}`));
    }
    catch (error) {
        spinner.fail(chalk_1.default.red(`Failed to block IP: ${ip}`));
        console.error(error);
    }
}
async function unblock(ip, options) {
    const spinner = (0, ora_1.default)(`Unblocking IP: ${ip}...`).start();
    try {
        const client = (0, redis_1.createClient)({ url: options.redis });
        await client.connect();
        await client.sRem('chaos:settings:blocked_ips', ip);
        await client.quit();
        spinner.succeed(chalk_1.default.green(`Unblocked IP: ${ip}`));
    }
    catch (error) {
        spinner.fail(chalk_1.default.red(`Failed to unblock IP: ${ip}`));
        console.error(error);
    }
}
async function listBlocked(options) {
    const spinner = (0, ora_1.default)('Fetching blocked IPs...').start();
    try {
        const client = (0, redis_1.createClient)({ url: options.redis });
        await client.connect();
        const ips = await client.sMembers('chaos:settings:blocked_ips');
        await client.quit();
        if (ips.length === 0) {
            spinner.info(chalk_1.default.yellow('No IPs are currently blocked.'));
            return;
        }
        spinner.stop();
        console.log(chalk_1.default.cyan('\n🚫 Blocked IPs:'));
        ips.forEach(ip => console.log(`   - ${ip}`));
        console.log('');
    }
    catch (error) {
        spinner.fail(chalk_1.default.red('Failed to fetch blocked IPs'));
        console.error(error);
    }
}
//# sourceMappingURL=block.js.map