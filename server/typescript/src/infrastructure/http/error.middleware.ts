import type { Request, Response, NextFunction } from 'express';
import type { ErrorResponseDTO } from './dtos.js';
import { AppError } from '../../shared/errors.js';

/**
 * Error handler middleware
 */
export function errorHandler(
  err: unknown,
  req: Request,
  res: Response,
  next: NextFunction
): void {
  console.error('Error:', err);

  // Handle known application errors
  if (err instanceof AppError) {
    const errorResponse: ErrorResponseDTO = {
      error: {
        code: err.code,
        message: err.message,
        details: err.details
      }
    };
    res.status(err.statusCode).json(errorResponse);
    return;
  }

  // Handle Zod validation errors
  if (err && typeof err === 'object' && 'issues' in err) {
    const details: Record<string, unknown> = {};
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    (err as any).issues.forEach((issue: { path: (string | number)[]; message: string }) => {
      const path = issue.path.join('.');
      details[path] = issue.message;
    });

    const errorResponse: ErrorResponseDTO = {
      error: {
        code: 'validation_error',
        message: 'Validation failed',
        details
      }
    };
    res.status(400).json(errorResponse);
    return;
  }

  // Handle unknown errors
  const errorResponse: ErrorResponseDTO = {
    error: {
      code: 'internal_error',
      message: 'An internal error occurred'
    }
  };
  res.status(500).json(errorResponse);
}

/**
 * Not found handler middleware
 */
export function notFoundHandler(req: Request, res: Response): void {
  const errorResponse: ErrorResponseDTO = {
    error: {
      code: 'not_found',
      message: `Route ${req.method} ${req.path} not found`
    }
  };
  res.status(404).json(errorResponse);
}
