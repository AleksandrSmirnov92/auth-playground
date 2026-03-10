"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.AuthUsecase = void 0;
const bcrypt_1 = __importDefault(require("bcrypt"));
const uuid_1 = require("uuid");
/**
 * AuthUsecase — бизнес-логика авторизации (UseCase слой).
 *
 * Здесь нет HTTP и Express: только правила приложения:
 * - register: проверка уникальности email + bcrypt hash + сохранение
 * - login: проверка email + bcrypt compare
 */
class AuthUsecase {
    userRepository;
    constructor(userRepository) {
        this.userRepository = userRepository;
    }
    async register(email, password) {
        const existing = await this.userRepository.getByEmail(email);
        if (existing) {
            throw new Error('user already exists');
        }
        const hashedPassword = await bcrypt_1.default.hash(password, 10);
        const user = {
            id: (0, uuid_1.v4)(),
            email,
            password: hashedPassword,
            createdAt: new Date(),
        };
        await this.userRepository.create(user);
        const { password: _, ...userWithoutPassword } = user;
        return userWithoutPassword;
    }
    async login(email, password) {
        const user = await this.userRepository.getByEmail(email);
        if (!user) {
            throw new Error('invalid email or password');
        }
        const valid = await bcrypt_1.default.compare(password, user.password);
        if (!valid) {
            throw new Error('invalid email or password');
        }
        return user;
    }
    async getUserById(id) {
        const user = await this.userRepository.getById(id);
        if (!user)
            return null;
        const { password: _, ...userWithoutPassword } = user;
        return userWithoutPassword;
    }
    async deleteUserById(id) {
        await this.userRepository.delete(id);
    }
}
exports.AuthUsecase = AuthUsecase;
