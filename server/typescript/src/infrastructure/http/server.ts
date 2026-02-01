import express from 'express';
import path from 'path';
import { fileURLToPath } from 'url';
import { createItemsRouter } from './routes.js';
import { errorHandler, notFoundHandler } from './error.middleware.js';
import { SQLiteItemRepository } from '../persistence/item.repository.impl.js';
import { ensureDatabaseInitialized } from '../persistence/database-setup.js';

/**
 * Create Express app
 */
export function createApp(): express.Express {
  const app = express();

  // Middleware
  app.use(express.json());
  app.use(express.urlencoded({ extended: true }));

  // Ensure database is initialized
  ensureDatabaseInitialized();

  // Create repository instance
  const itemRepository = new SQLiteItemRepository();

  const clientIndexPath = path.resolve(
    path.dirname(fileURLToPath(import.meta.url)),
    '../../../../..',
    'client',
    'index.html'
  );

  // Routes
  app.get('/', (req, res, next) => {
    res.sendFile(clientIndexPath, (err) => {
      if (err) {
        next(err);
      }
    });
  });
  app.use('/v1/items', createItemsRouter(express.Router(), itemRepository));

  // Health check
  app.get('/health', (req, res) => {
    res.json({ status: 'ok', timestamp: new Date().toISOString() });
  });

  // Error handlers (must be last)
  app.use(notFoundHandler);
  app.use(errorHandler);

  return app;
}

/**
 * Start server
 */
export async function startServer(port: number = 3000): Promise<void> {
  const app = createApp();

  return new Promise((resolve) => {
    const server = app.listen(port, () => {
      console.log(`Server listening on port ${port}`);
      resolve();
    });

    // Handle graceful shutdown
    process.on('SIGTERM', () => {
      console.log('SIGTERM received, shutting down gracefully...');
      server.close(() => {
        console.log('Server closed');
      });
    });
  });
}

// Start server if this file is run directly
if (import.meta.url === `file://${process.argv[1]}`) {
  const port = Number(process.env.PORT) || 3000;
  startServer(port).catch(console.error);
}
