import { User } from './user';

/**
 * UserRepository — контракт хранилища пользователей.
 *
 * UseCase зависит от интерфейса, а не от конкретной реализации,
 * поэтому можно заменить in-memory на Postgres/Redis без изменения бизнес-логики.
 */
export interface UserRepository {
  create(user: User): Promise<void>;
  getById(id: string): Promise<User | null>;
  getByEmail(email: string): Promise<User | null>;
  delete(id: string): Promise<void>;
}
