# LeBonPoint Marketplace Server (TypeScript)

TypeScript REST API server with SQLite database for the LeBonPoint marketplace application.

## Development

### Prerequisites

- Node.js 18+
- npm

### Installation

```bash
npm install
```

### Build

```bash
npm run build
```

### Run Development Server

```bash
npm run dev
```

The server will start on port 3000 (or the port specified in PORT environment variable).

### Run Tests

```bash
npm test
```

Run tests in watch mode:

```bash
npm run test:watch
```

## Environment Variables

- `PORT` - Server port (default: 3000)
- `NODE_ENV` - Environment (set to `test` for testing)
