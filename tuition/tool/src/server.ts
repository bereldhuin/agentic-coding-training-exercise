import express from 'express';
import bodyParser from 'body-parser';
import cors from 'cors';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';
import { exec } from 'child_process';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const app = express();
const PORT = process.env.PORT || 3000;
const DATA_FILE = path.join(__dirname, '../data/servers.json');

app.use(cors());
app.use(bodyParser.json());
app.use(express.static(path.join(__dirname, '../public')));

interface MCPServer {
  id: string;
  url: string;
  name: string;
  status: 'verified' | 'unverified' | 'error';
  lastChecked?: string;
}

const readData = (): MCPServer[] => {
  try {
    const data = fs.readFileSync(DATA_FILE, 'utf-8');
    return JSON.parse(data);
  } catch (e) {
    return [];
  }
};

const writeData = (data: MCPServer[]) => {
  fs.writeFileSync(DATA_FILE, JSON.stringify(data, null, 2));
};

app.get('/api/servers', (req, res) => {
  res.json(readData());
});

app.post('/api/register', (req, res) => {
  const { url, name } = req.body;
  if (!url) return res.status(400).json({ error: 'URL is required' });

  const servers = readData();
  const newServer: MCPServer = {
    id: Date.now().toString(),
    url,
    name: name || 'Unnamed Server',
    status: 'unverified'
  };
  servers.push(newServer);
  writeData(servers);
  res.status(201).json(newServer);
});

app.post('/api/verify/:id', async (req, res) => {
  const { id } = req.params;
  const servers = readData();
  const server = servers.find(s => s.id === id);

  if (!server) return res.status(404).json({ error: 'Server not found' });

  try {
    // Simple fetch check
    const response = await fetch(server.url, { method: 'HEAD' }).catch(() => null);
    server.status = response && response.ok ? 'verified' : 'error';
    server.lastChecked = new Date().toISOString();
    writeData(servers);
    res.json(server);
  } catch (error) {
    server.status = 'error';
    server.lastChecked = new Date().toISOString();
    writeData(servers);
    res.json(server);
  }
});

app.get('/api/ip-cmd', (req, res) => {
  // Command to get IP address on macOS/Linux
  const cmd = "ifconfig | grep 'inet ' | grep -v 127.0.0.1 | awk '{print $2}'";
  exec(cmd, (error, stdout, stderr) => {
    if (error) {
      return res.status(500).json({ error: error.message });
    }
    res.json({ output: stdout.trim(), command: cmd });
  });
});

app.listen(PORT, () => {
  console.log(`Server running at http://localhost:${PORT}`);
});
