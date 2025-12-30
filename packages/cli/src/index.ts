#!/usr/bin/env node

import { Command } from 'commander';
import chalk from 'chalk';
import { init } from './commands/init.js';
import { status } from './commands/status.js';
import { start } from './commands/start.js';
import { logs } from './commands/logs.js';
import { block, unblock, listBlocked } from './commands/block.js';

const program = new Command();

console.log(chalk.magenta(`
   _____ _                        _____                     
  / ____| |                      |  __ \\                    
 | |    | |__   __ _  ___  ___   | |__) | __ _____  ___   _ 
 | |    | '_ \\ / _\` |/ _ \\/ __|  |  ___/ '__/ _ \\ \\/ / | | |
 | |____| | | | (_| | (_) \\__ \\  | |   | | | (_) >  <| |_| |
  \\_____|_| |_|\\__,_|\\___/|___/  |_|   |_|  \\___/_/\\_\\\\__, |
                                                       __/ |
                                            👻        |___/ 
`));

program
    .name('chaos-proxy')
    .description('CLI for managing Chaos-Proxy - The Immortality Layer for APIs')
    .version('1.0.0');

program
    .command('init')
    .description('Initialize a new Chaos-Proxy project')
    .option('-d, --docker', 'Setup with Docker support', false)
    .action(init);

program
    .command('status')
    .description('Check the status of Chaos-Proxy services')
    .option('-r, --redis <url>', 'Redis connection URL', 'redis://localhost:6379')
    .action(status);

program
    .command('start')
    .description('Start Chaos-Proxy services')
    .option('--docker', 'Start using Docker Compose', false)
    .action(start);

program
    .command('logs')
    .description('View service logs (Docker only)')
    .option('-s, --service <name>', 'Filter by service (sentinel, brain, dashboard, redis)')
    .option('-f, --follow', 'Follow log output', true)
    .action(logs);

program
    .command('block')
    .description('Block an IP address')
    .argument('<ip>', 'IP address to block')
    .option('-r, --redis <url>', 'Redis connection URL', 'redis://localhost:6379')
    .action(block);

program
    .command('unblock')
    .description('Unblock an IP address')
    .argument('<ip>', 'IP address to unblock')
    .option('-r, --redis <url>', 'Redis connection URL', 'redis://localhost:6379')
    .action(unblock);

program
    .command('ls-blocked')
    .description('List all blocked IP addresses')
    .option('-r, --redis <url>', 'Redis connection URL', 'redis://localhost:6379')
    .action(listBlocked);

program.parse();
