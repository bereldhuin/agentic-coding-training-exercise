import type { IItemRepository } from '../../domain/repositories/item.repository.port.js';
import type { CreateItemData } from '../../domain/entities/item.entity.js';
import { createItemWithDefaults } from '../../domain/entities/item.entity.js';
import { ValidationService } from '../services/validation.service.js';
import type { Item } from '../../domain/entities/index.js';
import { Condition, Status } from '../../domain/entities/index.js';

/**
 * Create item use case
 */
export class CreateItemUseCase {
  constructor(private readonly itemRepository: IItemRepository) {}

  /**
   * Execute the use case
   */
  async execute(data: CreateItemData): Promise<Item> {
    // Validate input
    const validatedData = ValidationService.validateCreateItem(data);

    // Apply defaults with proper enum casting
    const itemData = createItemWithDefaults({
      ...validatedData,
      condition: validatedData.condition as Condition,
      status: validatedData.status as Status
    });

    // Create item in repository
    const createdItem = await this.itemRepository.create(itemData);

    return createdItem;
  }
}
