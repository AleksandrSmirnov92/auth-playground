import { User } from '../domain/user';
import { UserRepository } from '../domain/repository';
import bcrypt from 'bcrypt';
import { v4 as uuidv4 } from 'uuid';

export type PublicUser = {
  id: string;
  email: string;
  created_at: string; // ISO datetime string
};

/**
 * AuthUsecase — бизнес-логика авторизации (UseCase слой).
 *
 * Здесь нет HTTP и Express: только правила приложения:
 * - register: проверка уникальности email + bcrypt hash + сохранение
 * - login: проверка email + bcrypt compare
 */
export class AuthUsecase {
  constructor(private userRepository: UserRepository) {}

  async register(email: string, password: string): Promise<PublicUser> {
    const existing = await this.userRepository.getByEmail(email);
    if (existing) {
      throw new Error('user already exists');
    }

    const hashedPassword = await bcrypt.hash(password, 10);
    const user: User = {
      id: uuidv4(),
      email,
      password: hashedPassword,
      createdAt: new Date(),
    };

    await this.userRepository.create(user);
    return { id: user.id, email: user.email, created_at: user.createdAt.toISOString() };
  }

  async login(email: string, password: string): Promise<User> {
    const user = await this.userRepository.getByEmail(email);
    if (!user) {
      throw new Error('invalid email or password');
    }
    const valid = await bcrypt.compare(password, user.password);
    if (!valid) {
      throw new Error('invalid email or password');
    }
    return user;
  }

  async getUserById(id: string): Promise<PublicUser | null> {
    const user = await this.userRepository.getById(id);
    if (!user) return null;
    return { id: user.id, email: user.email, created_at: user.createdAt.toISOString() };
  }

  async deleteUserById(id: string): Promise<void> {
    await this.userRepository.delete(id);
  }

  // deleteByCredentials — проверяет email+пароль и удаляет пользователя.
  // Повторно использует login(), чтобы не дублировать логику bcrypt.
  async deleteByCredentials(email: string, password: string): Promise<void> {
    const user = await this.login(email, password);
    await this.deleteUserById(user.id);
  }
}
