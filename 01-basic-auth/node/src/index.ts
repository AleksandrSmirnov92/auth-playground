/**
 * Entry point (Node.js / TypeScript).
 *
 * Здесь мы:
 * - создаём зависимости: Repository → UseCase → Delivery
 * - регистрируем маршруты
 * - поднимаем HTTP сервер на порту 8080
 *
 * BasicAuth middleware удалён — авторизация выполняется через тело запроса.
 */
import express from 'express';
import fs from 'node:fs';
import path from 'node:path';
import { MemoryUserRepository } from './repository/memoryUserRepository';
import { AuthUsecase } from './usecase/authUsecase';
import { createAuthHandler } from './delivery/authHandler';

const userRepository = new MemoryUserRepository();
const authUsecase = new AuthUsecase(userRepository);
const authHandler = createAuthHandler(authUsecase);

const app = express();
app.use(express.json());

// GET /health — проверка работоспособности сервера
app.get('/health', (_req, res) => {
  res.json({ status: 'ok' });
});

// GET /openapi.json — отдаёт OpenAPI спецификацию из файла
app.get('/openapi.json', (_req, res) => {
  const specPath = path.join(process.cwd(), 'openapi.json');
  const data = fs.readFileSync(specPath, 'utf-8');
  res.type('application/json').send(data);
});

// GET /swagger — Swagger UI (HTML-страница)
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

// Публичные маршруты — не требуют авторизации
app.post('/api/v1/auth/register', authHandler.register);
app.post('/api/v1/auth/login', authHandler.login);
app.delete('/api/v1/auth/delete', authHandler.deleteByCredentials);

const PORT = 8080;
app.listen(PORT, () => {
  console.log(`Basic Auth server (Node/TypeScript) running on http://localhost:${PORT}`);
  console.log(`Swagger UI: http://localhost:${PORT}/swagger`);
});
