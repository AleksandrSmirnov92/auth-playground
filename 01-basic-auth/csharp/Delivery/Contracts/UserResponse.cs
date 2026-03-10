using System.Text.Json.Serialization;

namespace BasicAuth.Delivery.Contracts;

// PublicUserResponse — то, что мы отдаём наружу в JSON.
//
// Важно: имена полей фиксируем через JsonPropertyName, чтобы API было единообразным
// между языками (например, created_at как snake_case).
public record UserResponse(
    [property: JsonPropertyName("id")] string Id,
    [property: JsonPropertyName("email")] string Email,
    [property: JsonPropertyName("created_at")] DateTime CreatedAt
);

