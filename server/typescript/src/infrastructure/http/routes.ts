import type { Router, Request, Response, NextFunction } from 'express';
import type {
  CreateItemRequestDTO,
  UpdateItemRequestDTO,
  PatchItemRequestDTO
} from './dtos.js';
import { toItemResponseDTO } from './dtos.js';
import type { Item } from '../../domain/entities/item.entity.js';
import {
  CreateItemUseCase,
  GetItemUseCase,
  ListItemsUseCase,
  UpdateItemUseCase,
  PatchItemUseCase,
  DeleteItemUseCase
} from '../../application/use-cases/index.js';
import type { IItemRepository } from '../../domain/repositories/item.repository.port.js';
import { ValidationError } from '../../shared/errors.js';

/**
 * Create items routes
 */
export function createItemsRouter(router: Router, itemRepository: IItemRepository): Router {
  // Initialize use cases
  const createItemUseCase = new CreateItemUseCase(itemRepository);
  const getItemUseCase = new GetItemUseCase(itemRepository);
  const listItemsUseCase = new ListItemsUseCase(itemRepository);
  const updateItemUseCase = new UpdateItemUseCase(itemRepository);
  const patchItemUseCase = new PatchItemUseCase(itemRepository);
  const deleteItemUseCase = new DeleteItemUseCase(itemRepository);

  // GET /v1/items - List items with filters, sorting, pagination, and search
  router.get('/', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const result = await listItemsUseCase.execute(req.query);
      const items = result.items.map(item => toItemResponseDTO(item as Item));

      res.json({
        items,
        next_cursor: result.next_cursor
      });
    } catch (error) {
      next(error);
    }
  });

  // POST /v1/items - Create a new item
  router.post('/', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const data = req.body as CreateItemRequestDTO;
      const item = await createItemUseCase.execute(data as any);

      res.status(201).json(toItemResponseDTO(item));
    } catch (error) {
      next(error);
    }
  });

  // GET /v1/items/:id - Get a single item
  router.get('/:id', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const id = parseItemId(req.params.id);
      const item = await getItemUseCase.execute(id);

      res.json(toItemResponseDTO(item));
    } catch (error) {
      next(error);
    }
  });

  // PUT /v1/items/:id - Replace an item
  router.put('/:id', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const id = parseItemId(req.params.id);
      const data = req.body as UpdateItemRequestDTO;
      const item = await updateItemUseCase.execute(id, data as any);

      res.json(toItemResponseDTO(item));
    } catch (error) {
      next(error);
    }
  });

  // PATCH /v1/items/:id - Partially update an item
  router.patch('/:id', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const id = parseItemId(req.params.id);
      const data = req.body as PatchItemRequestDTO;
      const item = await patchItemUseCase.execute(id, data as any);

      res.json(toItemResponseDTO(item));
    } catch (error) {
      next(error);
    }
  });

  // DELETE /v1/items/:id - Delete an item
  router.delete('/:id', async (req: Request, res: Response, next: NextFunction) => {
    try {
      const id = parseItemId(req.params.id);
      await deleteItemUseCase.execute(id);

      res.status(204).send();
    } catch (error) {
      next(error);
    }
  });

  return router;
}

/**
 * Parse item ID from request params
 */
function parseItemId(value: string | string[]): number {
  const idValue = Array.isArray(value) ? value[0] : value;
  const id = Number(idValue);

  if (isNaN(id) || id <= 0) {
    throw new ValidationError('Invalid item ID', { id: 'ID must be a positive integer' });
  }

  return id;
}
