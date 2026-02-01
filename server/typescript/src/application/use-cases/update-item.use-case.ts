import type { IItemRepository } from '../../domain/repositories/item.repository.port.js';
import { NotFoundError, ValidationError } from '../../shared/errors.js';
import { ValidationService } from '../services/validation.service.js';
import type { ReplaceItemData } from '../../domain/entities/item.entity.js';
import type { Item } from '../../domain/entities/index.js';
import { Condition, Status } from '../../domain/entities/index.js';

/**
 * Update item (PUT - full replace) use case
 */
export class UpdateItemUseCase {
  constructor(private readonly itemRepository: IItemRepository) {}

  /**
   * Execute the use case
   */
  async execute(id: number, data: ReplaceItemData): Promise<Item> {
    // Validate input
    const validatedData = ValidationService.validateUpdateItem(data);

    // Check if item exists
    const exists = await this.itemRepository.exists(id);
    if (!exists) {
      throw new NotFoundError('Item', id);
    }

    // Cast enums and update item in repository
    const updateData = {
      ...validatedData,
      condition: validatedData.condition as Condition,
      status: validatedData.status as Status
    };

    const updatedItem = await this.itemRepository.update(id, updateData);

    if (!updatedItem) {
      throw new NotFoundError('Item', id);
    }

    return updatedItem;
  }
}
