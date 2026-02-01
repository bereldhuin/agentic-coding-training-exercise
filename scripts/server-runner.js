#!/usr/bin/env node

const path = require('path');
const readline = require('readline/promises');
const { stdin, stdout } = require('process');
const { DEFAULT_TIMEOUT_MS, servers } = require('./lib/servers');
const {
  buildIfNeeded,
  logError,
  logInfo,
  startServer,
  stopProcess,
  verifyServerHealth,
  waitForHealth,
} = require('./lib/runner');
const { runApiVerification, testCategories } = require('./lib/api-tests');

const REPO_ROOT = path.resolve(__dirname, '..');
process.chdir(REPO_ROOT);

function parseArgs() {
  const args = process.argv.slice(2);
  const config = {
    mode: null,
    servers: servers.map((server) => server.id),
    tests: Object.keys(testCategories),
    timeoutMs: DEFAULT_TIMEOUT_MS,
    verbose: false,
    tap: false,
    output: null,
    serverId: null,
  };

  for (let i = 0; i < args.length; i++) {
    const arg = args[i];

    switch (arg) {
      case '--mode':
        i++;
        config.mode = args[i];
        break;
      case '--servers':
        i++;
        config.servers = args[i].split(',').map((s) => s.trim());
        break;
      case '--server':
        i++;
        config.serverId = args[i].trim();
        break;
      case '--tests':
        i++;
        config.tests = args[i].split(',').map((t) => t.trim());
        break;
      case '--timeout':
        i++;
        config.timeoutMs = Number(args[i]) * 1000;
        break;
      case '--verbose':
      case '-v':
        config.verbose = true;
        break;
      case '--tap':
        config.tap = true;
        break;
      case '--output':
        i++;
        config.output = args[i];
        break;
      case '--verify-api':
        config.mode = 'verify-api';
        break;
      case '--verify-health':
        config.mode = 'verify-health';
        break;
      case '--help':
      case '-h':
        config.mode = 'help';
        break;
      default:
        logError(`Unknown option: ${arg}`);
        config.mode = 'help';
        break;
    }
  }

  return config;
}

function printHelp() {
  console.log(`
Server Runner

Usage:
  node scripts/server-runner.js
  node scripts/server-runner.js --verify-health
  node scripts/server-runner.js --verify-api [--servers typescript,python] [--tests health,list]

Options:
  --mode <start|verify-health|verify-api>
  --server <id>           Start a single server (start mode)
  --servers <list>        Comma-separated server ids (verify modes)
  --tests <list>          Comma-separated test categories (verify-api)
  --timeout <seconds>     Timeout per server (default: 30)
  --verbose               Verbose output (verify-api)
  --tap                   TAP output (verify-api)
  --output <file>         JSON report output (verify-api)
  --help                  Show this help message

Available servers: ${servers.map((server) => server.id).join(', ')}
Available test categories: ${Object.keys(testCategories).join(', ')}
`);
}

function resolveServers(ids) {
  const selected = servers.filter((server) => ids.includes(server.id));
  if (selected.length === 0) {
    throw new Error('No matching servers selected.');
  }
  return selected;
}

async function verifyAllServers(timeoutMs) {
  const results = [];

  for (const server of servers) {
    try {
      await verifyServerHealth(server, timeoutMs);
      results.push({ server: server.id, ok: true });
    } catch (error) {
      results.push({ server: server.id, ok: false, error: error.message });
    }
  }

  const failed = results.filter((result) => !result.ok);
  if (failed.length > 0) {
    logError('\nVerification failed:');
    for (const result of failed) {
      logError(`- ${result.server}: ${result.error}`);
    }
    process.exit(1);
  }

  logInfo('\nAll servers verified successfully.');
}

async function promptMode() {
  const rl = readline.createInterface({ input: stdin, output: stdout });
  const mode = await rl.question(
    'Choose an action:\n  1) Start a server\n  2) Verify health (all servers)\n  3) Verify API (all servers)\n\nEnter 1, 2, or 3: '
  );

  if (mode.trim() === '2') {
    rl.close();
    return { mode: 'verify-health' };
  }

  if (mode.trim() === '3') {
    rl.close();
    return { mode: 'verify-api' };
  }

  const choices = servers
    .map((server, index) => `  ${index + 1}) ${server.label}`)
    .join('\n');

  const answer = await rl.question(`\nSelect a server:\n${choices}\n\nEnter a number: `);
  rl.close();

  const choiceIndex = Number(answer.trim()) - 1;
  if (Number.isNaN(choiceIndex) || choiceIndex < 0 || choiceIndex >= servers.length) {
    throw new Error('Invalid server selection.');
  }

  return { mode: 'start', server: servers[choiceIndex] };
}

async function startInteractiveServer(server, timeoutMs) {
  await buildIfNeeded(server);
  const child = await startServer(server);
  const healthy = await waitForHealth(server, timeoutMs);

  if (!healthy) {
    stopProcess(child, server);
    throw new Error(`${server.id} failed health check within ${timeoutMs / 1000}s`);
  }

  logInfo(`${server.label} is running. Press Ctrl+C to stop.`);

  const handleExit = async () => {
    await stopProcess(child, server);
    process.exit(0);
  };

  process.on('SIGINT', () => {
    void handleExit();
  });
  process.on('SIGTERM', () => {
    void handleExit();
  });

  await new Promise(() => {});
}

async function run() {
  try {
    const config = parseArgs();

    if (config.mode === 'help') {
      printHelp();
      return;
    }

    if (config.mode === 'verify-api') {
      const selected = resolveServers(config.servers);
      await runApiVerification(selected, {
        tests: config.tests,
        timeoutMs: config.timeoutMs,
        verbose: config.verbose,
        tap: config.tap,
        output: config.output,
      });
      return;
    }

    if (config.mode === 'verify-health') {
      const selected = resolveServers(config.servers);
      for (const server of selected) {
        await verifyServerHealth(server, config.timeoutMs);
      }
      logInfo('\nAll selected servers verified successfully.');
      return;
    }

    if (config.mode === 'start' && config.serverId) {
      const server = resolveServers([config.serverId])[0];
      await startInteractiveServer(server, config.timeoutMs);
      return;
    }

    const prompt = await promptMode();

    if (prompt.mode === 'verify-health') {
      await verifyAllServers(config.timeoutMs);
      return;
    }

    if (prompt.mode === 'verify-api') {
      await runApiVerification(servers, {
        tests: Object.keys(testCategories),
        timeoutMs: config.timeoutMs,
        verbose: false,
        tap: false,
        output: null,
      });
      return;
    }

    await startInteractiveServer(prompt.server, config.timeoutMs);
  } catch (error) {
    logError(`\nError: ${error.message}`);
    process.exit(1);
  }
}

run();
