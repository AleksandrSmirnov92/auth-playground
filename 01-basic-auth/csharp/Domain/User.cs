namespace BasicAuth.Domain;

// User — доменная сущность пользователя.
//
// Зачем она нужна:
// - Repository хранит/ищет пользователей
// - UseCase создаёт пользователя при регистрации и проверяет пароль при логине
// - Delivery отдаёт данные пользователя в JSON (пароль наружу не возвращаем)
public class User
{
    public required string Id { get; set; }
    public required string Email { get; set; }
    public required string Password { get; set; }
    public DateTime CreatedAt { get; set; }
}
