/**
 * User — доменная сущность пользователя.
 *
 * Зачем она нужна:
 * - Repository хранит/ищет пользователей
 * - UseCase создаёт пользователя при регистрации и проверяет пароль при логине
 * - Delivery отдаёт данные пользователя в JSON (в ответах пароль не возвращаем)
 */
export interface User {
  id: string;
  email: string;
  password: string;
  createdAt: Date;
}
