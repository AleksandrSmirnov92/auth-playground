import { Request, Response } from 'express';
import { AuthUsecase } from '../usecase/authUsecase';
import { USER_ID_KEY } from './middleware/basicAuth';

export function createAuthHandler(authUsecase: AuthUsecase) {
  return {
    async register(req: Request, res: Response) {
      const { email, password } = req.body;
      if (!email || !password) {
        res.status(400).json({ error: 'email and password required' });
        return;
      }
      try {
        const user = await authUsecase.register(email, password);
        res.status(201).json(user);
      } catch (err: any) {
        if (err.message === 'user already exists') {
          res.status(409).json({ error: err.message });
        } else {
          res.status(500).json({ error: err.message });
        }
      }
    },

    async me(req: Request, res: Response) {
      const userId = (req as any)[USER_ID_KEY];
      if (!userId) {
        res.status(401).json({ error: 'unauthorized' });
        return;
      }
      const user = await authUsecase.getUserById(userId);
      if (!user) {
        res.status(404).json({ error: 'user not found' });
        return;
      }
      res.json(user);
    },

    async deleteUser(req: Request, res: Response) {
      const userId = (req as any)[USER_ID_KEY];
      if (!userId) {
        res.status(401).json({ error: 'unauthorized' });
        return;
      }
      await authUsecase.deleteUserById(userId);
      res.status(204).send();
    },
  };
}
