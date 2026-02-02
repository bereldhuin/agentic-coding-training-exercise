const path = require('path');

const DEFAULT_TIMEOUT_MS = 30000;
const HEALTH_PATH = '/health';

const servers = [
  {
    id: 'typescript',
    label: 'TypeScript (Express)',
    cwd: 'server/typescript',
    port: 3000,
    healthPath: HEALTH_PATH,
    build: { command: 'npm', args: ['run', 'build'] },
    start: { command: 'node', args: ['dist/infrastructure/http/server.js'] },
    env: {},
  },
  {
    id: 'go',
    label: 'Go (Gin)',
    cwd: 'server/go',
    port: 8081,
    healthPath: HEALTH_PATH,
    build: { command: 'go', args: ['build', '-tags', 'fts5', '-o', 'bin/server', 'cmd/server/main.go'] },
    start: { command: './bin/server', args: [] },
    env: {},
  },
  {
    id: 'python',
    label: 'Python (FastAPI)',
    cwd: 'server/python',
    port: 8000,
    healthPath: HEALTH_PATH,
    build: null,
    start: {
      command: 'poetry',
      args: ['run', 'uvicorn', 'src.main:app', '--host', '0.0.0.0', '--port', '8000'],
    },
    env: { HOST: '0.0.0.0' },
  },
  {
    id: 'swift',
    label: 'Swift (Vapor)',
    cwd: 'server/swift',
    port: 9000,
    healthPath: HEALTH_PATH,
    build: { command: 'swift', args: ['build'] },
    start: { command: '.build/debug/App', args: [] },
    env: {},
  },
  {
    id: 'kotlin',
    label: 'Kotlin (Ktor)',
    cwd: 'server/kotlin',
    port: 8080,
    healthPath: HEALTH_PATH,
    build: { command: './gradlew', args: ['installDist'] },
    start: { command: './build/install/lebonpoint-kotlin-server/bin/lebonpoint-kotlin-server', args: [] },
    env: {},
    needsJavaHome: true,
  },
];

module.exports = {
  DEFAULT_TIMEOUT_MS,
  HEALTH_PATH,
  servers,
};
