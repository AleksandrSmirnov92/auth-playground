using System.Text.Json.Serialization;

namespace BasicAuth.Delivery.Contracts;

// AuthResponse — ответ на успешный register или login.
//
// Возвращает сообщение (на русском) и публичные данные пользователя.
// Пример: { "message": "Пользователь успешно зарегистрирован", "user": { ... } }
public record AuthResponse(
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("user")]    UserResponse User
);
