/**
 * Entry point (Node.js / TypeScript).
 *
 * Здесь мы:
 * - создаём зависимости (Repository → UseCase → Delivery)
 * - подключаем middleware Basic Auth к защищённым роутам
 * - поднимаем HTTP сервер на порту 8080
 */
import express from 'express';
import fs from 'node:fs';
import path from 'node:path';
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

app.get('/openapi.json', (_req, res) => {
  const specPath = path.join(process.cwd(), 'openapi.json');
  const data = fs.readFileSync(specPath, 'utf-8');
  res.type('application/json').send(data);
});

app.get('/swagger', (_req, res) => {
  res.type('text/html').send(`<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Swagger UI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.onload = () => {
        SwaggerUIBundle({ url: '/openapi.json', dom_id: '#swagger-ui' });
      };
    </script>
  </body>
</html>`);
});

app.post('/api/v1/auth/register', authHandler.register);
app.get('/api/v1/auth/me', basicAuth, authHandler.me);
app.delete('/api/v1/auth/me', basicAuth, authHandler.deleteUser);

const PORT = 8080;
app.listen(PORT, () => {
  console.log(`Basic Auth server (Node/TypeScript) running on http://localhost:${PORT}`);
});
