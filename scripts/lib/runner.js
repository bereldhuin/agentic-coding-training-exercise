const { spawn } = require('child_process');

const ROOT_CWD = process.cwd();

function logInfo(message) {
  console.log(message);
}

function logError(message) {
  console.error(message);
}

function resolveEnv(server) {
  const env = {
    ...process.env,
    PORT: String(server.port),
    ...server.env,
  };

  if (server.needsJavaHome && process.env.JAVA_HOME) {
    env.JAVA_HOME = process.env.JAVA_HOME;
    if (!env.PATH?.includes(`${env.JAVA_HOME}/bin`)) {
      env.PATH = `${env.JAVA_HOME}/bin:${env.PATH || ''}`;
    }
  }

  return env;
}

function spawnCommand(command, args, options = {}) {
  return spawn(command, args, {
    stdio: options.stdio || 'inherit',
    cwd: options.cwd || ROOT_CWD,
    env: options.env || process.env,
    detached: Boolean(options.detached),
  });
}

async function runCommand(command, args, options = {}) {
  return new Promise((resolve, reject) => {
    const child = spawnCommand(command, args, options);
    child.on('error', reject);
    child.on('exit', (code) => {
      if (code === 0) {
        resolve();
        return;
      }
      reject(new Error(`${command} exited with code ${code}`));
    });
  });
}

async function waitForHealth(server, timeoutMs) {
  const url = `http://127.0.0.1:${server.port}${server.healthPath || '/health'}`;
  const end = Date.now() + timeoutMs;

  while (Date.now() < end) {
    try {
      const controller = new AbortController();
      const timeout = setTimeout(() => controller.abort(), 2000);
      const response = await fetch(url, { signal: controller.signal });
      clearTimeout(timeout);
      if (response.ok) {
        return true;
      }
    } catch (error) {
      // Keep retrying until timeout.
    }
    await new Promise((resolve) => setTimeout(resolve, 500));
  }

  return false;
}

async function stopProcess(child, server) {
  if (!child || child.killed) {
    return;
  }

  const signalOrder = server?.id === 'swift'
    ? ['SIGINT', 'SIGTERM', 'SIGKILL']
    : ['SIGTERM', 'SIGKILL'];

  for (const signal of signalOrder) {
    try {
      if (child.pid) {
        process.kill(-child.pid, signal);
      } else {
        child.kill(signal);
      }
    } catch (error) {
      return;
    }

    const exited = await new Promise((resolve) => {
      const timeout = setTimeout(() => resolve(false), 4000);
      child.once('exit', () => {
        clearTimeout(timeout);
        resolve(true);
      });
    });

    if (exited) {
      return;
    }
  }
}

async function buildIfNeeded(server) {
  if (!server.build) {
    return;
  }

  logInfo(`\nBuilding ${server.label}...`);
  try {
    await runCommand(server.build.command, server.build.args, {
      cwd: server.cwd,
      env: resolveEnv(server),
    });
  } catch (error) {
    throw new Error(`${server.id} build failed: ${error.message}`);
  }
}

async function startServer(server) {
  logInfo(`\nStarting ${server.label} on 0.0.0.0:${server.port}...`);
  const child = spawnCommand(server.start.command, server.start.args, {
    cwd: server.cwd,
    env: resolveEnv(server),
    detached: true,
  });

  child.on('error', (error) => {
    logError(`Failed to start ${server.label}: ${error.message}`);
  });

  return child;
}

async function verifyServerHealth(server, timeoutMs) {
  logInfo(`\nVerifying ${server.label}...`);

  await buildIfNeeded(server);
  const child = await startServer(server);
  const healthy = await waitForHealth(server, timeoutMs);

  if (!healthy) {
    await stopProcess(child, server);
    throw new Error(`${server.id} failed health check within ${timeoutMs / 1000}s`);
  }

  logInfo(`${server.label} is healthy.`);
  await stopProcess(child, server);
}

module.exports = {
  buildIfNeeded,
  logError,
  logInfo,
  resolveEnv,
  runCommand,
  spawnCommand,
  startServer,
  stopProcess,
  verifyServerHealth,
  waitForHealth,
};
