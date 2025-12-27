#!/usr/bin/env node
"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const commander_1 = require("commander");
const chalk_1 = __importDefault(require("chalk"));
const init_js_1 = require("./commands/init.js");
const status_js_1 = require("./commands/status.js");
const start_js_1 = require("./commands/start.js");
const program = new commander_1.Command();
console.log(chalk_1.default.magenta(`
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
    .action(init_js_1.init);
program
    .command('status')
    .description('Check the status of Chaos-Proxy services')
    .option('-r, --redis <url>', 'Redis connection URL', 'redis://localhost:6379')
    .action(status_js_1.status);
program
    .command('start')
    .description('Start Chaos-Proxy services')
    .option('--docker', 'Start using Docker Compose', false)
    .action(start_js_1.start);
program.parse();
//# sourceMappingURL=index.js.map