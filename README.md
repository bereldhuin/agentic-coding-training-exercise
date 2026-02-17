# Lebonpoint App

A multi-language classifieds platform implementation demonstrating consistent API behavior across different back-end technologies.

## Project Overview

Lebonpoint is a proof-of-concept application that implements a unified items API across five different programming languages and frameworks. The project focuses on API consistency, cross-server verification, and shared infrastructure.

### Supported Backend Implementations
- **TypeScript**: Node.js implementation
- **Python**: FastAPI implementation
- **Go**: Gin implementation
- **Kotlin**: Ktor implementation
- **Swift**: Vapor implementation

## Project Structure

```text
├── client/          # Frontend Single Page Application
├── scripts/         # Shared utility scripts (DB init, seeding, runner)
├── server/          # Backend implementations
│   ├── go/          # Go implementation
│   ├── kotlin/      # Kotlin implementation
│   ├── python/      # Python implementation
│   ├── swift/       # Swift implementation
│   └── typescript/  # TypeScript implementation
└── tuition/         # Server management and monitoring dashboard
```

## Getting Started

### Prerequisites
- **Node.js 18+**
- **PM2** (global): `npm install -g pm2`
- Language-specific runtimes (Go, Python, Kotlin/Java, Swift) for their respective implementations.

### Initial Setup

Install root dependencies:
```bash
npm install
```

### Running Servers
You can use the interactive server runner to start specific implementations or verify them all:
```bash
npm run server:run
```

### API Verification
To run the cross-server verification suite and ensure all implementations are compatible:
```bash
npm run verify
```

## Database Scripts

### Initialize Database

```bash
npm run db:init
```

### Seed Database

```bash
npm run db:seed
```

### Verify Database

```bash
npm run db:verify
```

## Documentation
More detailed documentation can be found in the following locations:
- [Verification Guide](./server/README-verification.md)
- [Project Specs](./openspec/specs/)

## License
MIT
