import { User } from './user';

export interface UserRepository {
  create(user: User): Promise<void>;
  getById(id: string): Promise<User | null>;
  getByEmail(email: string): Promise<User | null>;
  delete(id: string): Promise<void>;
}
