"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MemoryUserRepository = void 0;
/**
 * In-memory реализация UserRepository.
 *
 * Важно:
 * - данные живут только в памяти процесса (после перезапуска сервера исчезнут)
 * - getByEmail делает перебор пользователей (O(n)) — для учебного примера ок
 */
class MemoryUserRepository {
    users = new Map();
    async create(user) {
        this.users.set(user.id, user);
    }
    async getById(id) {
        return this.users.get(id) ?? null;
    }
    async getByEmail(email) {
        for (const user of this.users.values()) {
            if (user.email === email)
                return user;
        }
        return null;
    }
    async delete(id) {
        this.users.delete(id);
    }
}
exports.MemoryUserRepository = MemoryUserRepository;
