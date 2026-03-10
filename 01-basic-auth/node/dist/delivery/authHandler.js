"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.createAuthHandler = createAuthHandler;
const basicAuth_1 = require("./middleware/basicAuth");
function createAuthHandler(authUsecase) {
    return {
        async register(req, res) {
            const { email, password } = req.body;
            if (!email || !password) {
                res.status(400).json({ error: 'email and password required' });
                return;
            }
            try {
                const user = await authUsecase.register(email, password);
                res.status(201).json(user);
            }
            catch (err) {
                if (err.message === 'user already exists') {
                    res.status(409).json({ error: err.message });
                }
                else {
                    res.status(500).json({ error: err.message });
                }
            }
        },
        async me(req, res) {
            const userId = req[basicAuth_1.USER_ID_KEY];
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
        async deleteUser(req, res) {
            const userId = req[basicAuth_1.USER_ID_KEY];
            if (!userId) {
                res.status(401).json({ error: 'unauthorized' });
                return;
            }
            await authUsecase.deleteUserById(userId);
            res.status(204).send();
        },
    };
}
