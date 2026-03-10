import express from 'express';
import { MemoryUserRepository } from './repository/memoryUserRepository';
import { AuthUsecase } from './usecase/authUsecase';
import { createAuthHandler } from './delivery/authHandler';
import { basicAuthMiddleware } from './delivery/middleware/basicAuth';

const userRepository = new MemoryUserRepository();
const authUsecase = new AuthUsecase(userRepository);
const authHandler = createAuthHandler(authUsecase);
const basicAuth = basicAuthMiddleware(authUsecase);

const app = express();
app.use(express.json());

app.get('/health', (_req, res) => {
  res.json({ status: 'ok' });
});

app.post('/api/v1/auth/register', authHandler.register);
app.get('/api/v1/auth/me', basicAuth, authHandler.me);
app.delete('/api/v1/auth/me', basicAuth, authHandler.deleteUser);

const PORT = 8080;
app.listen(PORT, () => {
  console.log(`Basic Auth server (Node/TypeScript) running on http://localhost:${PORT}`);
});
