import type { IItemRepository } from '../../domain/repositories/item.repository.port.js';
import { NotFoundError } from '../../shared/errors.js';
import type { Item } from '../../domain/entities/index.js';

/**
 * Get item use case
 */
export class GetItemUseCase {
  constructor(private readonly itemRepository: IItemRepository) {}

  /**
   * Execute the use case
   */
  async execute(id: number): Promise<Item> {
    const item = await this.itemRepository.findById(id);

    if (!item) {
      throw new NotFoundError('Item', id);
    }

    return item;
  }
}
