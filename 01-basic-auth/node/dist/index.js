"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
/**
 * Entry point (Node.js / TypeScript).
 *
 * Здесь мы:
 * - создаём зависимости (Repository → UseCase → Delivery)
 * - подключаем middleware Basic Auth к защищённым роутам
 * - поднимаем HTTP сервер на порту 8080
 */
const express_1 = __importDefault(require("express"));
const node_fs_1 = __importDefault(require("node:fs"));
const node_path_1 = __importDefault(require("node:path"));
const memoryUserRepository_1 = require("./repository/memoryUserRepository");
const authUsecase_1 = require("./usecase/authUsecase");
const authHandler_1 = require("./delivery/authHandler");
const basicAuth_1 = require("./delivery/middleware/basicAuth");
const userRepository = new memoryUserRepository_1.MemoryUserRepository();
const authUsecase = new authUsecase_1.AuthUsecase(userRepository);
const authHandler = (0, authHandler_1.createAuthHandler)(authUsecase);
const basicAuth = (0, basicAuth_1.basicAuthMiddleware)(authUsecase);
const app = (0, express_1.default)();
app.use(express_1.default.json());
app.get('/health', (_req, res) => {
    res.json({ status: 'ok' });
});
app.get('/openapi.json', (_req, res) => {
    const specPath = node_path_1.default.join(process.cwd(), 'openapi.json');
    const data = node_fs_1.default.readFileSync(specPath, 'utf-8');
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
