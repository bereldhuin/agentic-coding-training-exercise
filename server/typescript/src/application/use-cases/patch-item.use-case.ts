import type { IItemRepository } from '../../domain/repositories/item.repository.port.js';
import { NotFoundError } from '../../shared/errors.js';
import { ValidationService } from '../services/validation.service.js';
import type { UpdateItemData } from '../../domain/entities/item.entity.js';
import type { Item } from '../../domain/entities/index.js';
import { Condition, Status } from '../../domain/entities/index.js';

/**
 * Patch item (partial update) use case
 */
export class PatchItemUseCase {
  constructor(private readonly itemRepository: IItemRepository) {}

  /**
   * Execute the use case
   */
  async execute(id: number, data: Partial<UpdateItemData>): Promise<Item> {
    // Validate input
    const validatedData = ValidationService.validatePatchItem(data);

    // Check if item exists
    const exists = await this.itemRepository.exists(id);
    if (!exists) {
      throw new NotFoundError('Item', id);
    }

    // Cast enums and patch item in repository
    const patchData: Partial<UpdateItemData> = {};
    if (validatedData.condition !== undefined) {
      patchData.condition = validatedData.condition as Condition;
    }
    if (validatedData.status !== undefined) {
      patchData.status = validatedData.status as Status;
    }
    // Add other fields
    Object.assign(patchData, validatedData);

    const updatedItem = await this.itemRepository.patch(id, patchData);

    if (!updatedItem) {
      throw new NotFoundError('Item', id);
    }

    return updatedItem;
  }
}
