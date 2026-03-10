import { Request, Response } from 'express';
import { AuthUsecase } from '../usecase/authUsecase';

/**
 * HTTP handlers (Delivery слой) для Express.
 *
 * Контракт:
 *   POST   /api/v1/auth/register  body: { email, password } → 201 { message, user }
 *   POST   /api/v1/auth/login     body: { email, password } → 200 { message, user }
 *   DELETE /api/v1/auth/delete    body: { email, password } → 200 { message }
 *
 * Handler'ы читают тело запроса, вызывают use case и формируют HTTP-ответ.
 * Basic Auth middleware больше не используется.
 */
export function createAuthHandler(authUsecase: AuthUsecase) {
  return {
    // POST /api/v1/auth/register
    // Создаёт нового пользователя. Возвращает 201 и данные пользователя.
    async register(req: Request, res: Response) {
      const { email, password } = req.body;
      if (!email || !password) {
        res.status(400).json({ error: 'email and password required' });
        return;
      }
      try {
        const user = await authUsecase.register(email, password);
        res.status(201).json({
          message: 'Пользователь успешно зарегистрирован',
          user,
        });
      } catch (err: any) {
        if (err.message === 'user already exists') {
          res.status(409).json({ error: err.message });
        } else {
          res.status(500).json({ error: err.message });
        }
      }
    },

    // POST /api/v1/auth/login
    // Проверяет email + пароль. Возвращает 200 и данные пользователя.
    async login(req: Request, res: Response) {
      const { email, password } = req.body;
      if (!email || !password) {
        res.status(400).json({ error: 'email and password required' });
        return;
      }
      try {
        const user = await authUsecase.login(email, password);
        const publicUser = {
          id: user.id,
          email: user.email,
          created_at: user.createdAt.toISOString(),
        };
        res.status(200).json({
          message: 'Добро пожаловать!',
          user: publicUser,
        });
      } catch (err: any) {
        if (err.message === 'invalid email or password') {
          res.status(401).json({ error: err.message });
        } else {
          res.status(500).json({ error: err.message });
        }
      }
    },

    // DELETE /api/v1/auth/delete
    // Принимает email + пароль в теле — проверяет и удаляет аккаунт.
    async deleteByCredentials(req: Request, res: Response) {
      const { email, password } = req.body;
      if (!email || !password) {
        res.status(400).json({ error: 'email and password required' });
        return;
      }
      try {
        await authUsecase.deleteByCredentials(email, password);
        res.status(200).json({ message: 'Пользователь успешно удалён' });
      } catch (err: any) {
        if (err.message === 'invalid email or password') {
          res.status(401).json({ error: err.message });
        } else {
          res.status(500).json({ error: err.message });
        }
      }
    },
  };
}
