import { Request, Response, NextFunction } from 'express';
import { AuthUsecase } from '../../usecase/authUsecase';

export const USER_ID_KEY = 'userId';

export function basicAuthMiddleware(authUsecase: AuthUsecase) {
  return (req: Request, res: Response, next: NextFunction) => {
    const authHeader = req.headers.authorization;
    if (!authHeader || !authHeader.startsWith('Basic ')) {
      res.setHeader('WWW-Authenticate', 'Basic realm="Restricted"');
      res.status(401).json({ error: 'Authorization required' });
      return;
    }

    const base64 = authHeader.slice(6);
    const decoded = Buffer.from(base64, 'base64').toString('utf-8');
    const [email, password] = decoded.split(':');
    if (!email || !password) {
      res.setHeader('WWW-Authenticate', 'Basic realm="Restricted"');
      res.status(401).json({ error: 'Invalid credentials' });
      return;
    }

    authUsecase
      .login(email, password)
      .then((user) => {
        (req as any)[USER_ID_KEY] = user.id;
        next();
      })
      .catch(() => {
        res.setHeader('WWW-Authenticate', 'Basic realm="Restricted"');
        res.status(401).json({ error: 'Invalid credentials' });
      });
  };
}
