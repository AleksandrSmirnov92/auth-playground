"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = __importDefault(require("express"));
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
app.post('/api/v1/auth/register', authHandler.register);
app.get('/api/v1/auth/me', basicAuth, authHandler.me);
app.delete('/api/v1/auth/me', basicAuth, authHandler.deleteUser);
const PORT = 8080;
app.listen(PORT, () => {
    console.log(`Basic Auth server (Node/TypeScript) running on http://localhost:${PORT}`);
});
