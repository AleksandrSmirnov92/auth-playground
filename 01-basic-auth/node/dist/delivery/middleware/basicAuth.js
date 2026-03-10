"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.USER_ID_KEY = void 0;
exports.basicAuthMiddleware = basicAuthMiddleware;
exports.USER_ID_KEY = 'userId';
function basicAuthMiddleware(authUsecase) {
    return (req, res, next) => {
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
            req[exports.USER_ID_KEY] = user.id;
            next();
        })
            .catch(() => {
            res.setHeader('WWW-Authenticate', 'Basic realm="Restricted"');
            res.status(401).json({ error: 'Invalid credentials' });
        });
    };
}
