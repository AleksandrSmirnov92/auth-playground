import { User } from '../domain/user';
import { UserRepository } from '../domain/repository';

export class MemoryUserRepository implements UserRepository {
  private users = new Map<string, User>();

  async create(user: User): Promise<void> {
    this.users.set(user.id, user);
  }

  async getById(id: string): Promise<User | null> {
    return this.users.get(id) ?? null;
  }

  async getByEmail(email: string): Promise<User | null> {
    for (const user of this.users.values()) {
      if (user.email === email) return user;
    }
    return null;
  }

  async delete(id: string): Promise<void> {
    this.users.delete(id);
  }
}
