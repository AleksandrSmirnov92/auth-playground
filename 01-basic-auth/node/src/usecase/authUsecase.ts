import { User } from '../domain/user';
import { UserRepository } from '../domain/repository';
import bcrypt from 'bcrypt';
import { v4 as uuidv4 } from 'uuid';

/**
 * AuthUsecase — бизнес-логика авторизации (UseCase слой).
 *
 * Здесь нет HTTP и Express: только правила приложения:
 * - register: проверка уникальности email + bcrypt hash + сохранение
 * - login: проверка email + bcrypt compare
 */
export class AuthUsecase {
  constructor(private userRepository: UserRepository) {}

  async register(email: string, password: string): Promise<Omit<User, 'password'>> {
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
    const { password: _, ...userWithoutPassword } = user;
    return userWithoutPassword;
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

  async getUserById(id: string): Promise<Omit<User, 'password'> | null> {
    const user = await this.userRepository.getById(id);
    if (!user) return null;
    const { password: _, ...userWithoutPassword } = user;
    return userWithoutPassword;
  }

  async deleteUserById(id: string): Promise<void> {
    await this.userRepository.delete(id);
  }
}
