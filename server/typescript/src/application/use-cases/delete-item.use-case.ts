import type { IItemRepository } from '../../domain/repositories/item.repository.port.js';
import { NotFoundError } from '../../shared/errors.js';

/**
 * Delete item use case
 */
export class DeleteItemUseCase {
  constructor(private readonly itemRepository: IItemRepository) {}

  /**
   * Execute the use case
   */
  async execute(id: number): Promise<void> {
    // Check if item exists
    const exists = await this.itemRepository.exists(id);
    if (!exists) {
      throw new NotFoundError('Item', id);
    }

    // Delete item from repository
    const deleted = await this.itemRepository.delete(id);

    if (!deleted) {
      throw new NotFoundError('Item', id);
    }
  }
}
