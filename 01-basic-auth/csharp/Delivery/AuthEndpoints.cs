using System.Text.Json.Serialization;
using BasicAuth.Delivery.Contracts;
using BasicAuth.UseCase;

namespace BasicAuth.Delivery;

// AuthEndpoints — маппинг URL → обработчик (Delivery / Minimal API слой).
//
// Контракт API (все body — JSON):
//   POST   /api/v1/auth/register  { email, password }  → 201 AuthResponse
//   POST   /api/v1/auth/login     { email, password }  → 200 AuthResponse
//   DELETE /api/v1/auth/delete    { email, password }  → 200 { message }
//
// Basic Auth middleware больше не используется — каждый запрос самодостаточен.
public static class AuthEndpoints
{
    public static void MapAuthEndpoints(this IEndpointRouteBuilder app, AuthUsecase authUsecase)
    {
        app.MapGet("/health", () => Results.Json(new { status = "ok" }));

        // POST /api/v1/auth/register
        // Принимает email + password, хеширует пароль и сохраняет пользователя.
        // 201 — пользователь создан; 409 — уже существует; 400 — пустые поля.
        app.MapPost("/api/v1/auth/register", async (RegisterRequest req) =>
        {
            if (string.IsNullOrEmpty(req.Email) || string.IsNullOrEmpty(req.Password))
                return Results.BadRequest(new { error = "email and password required" });

            try
            {
                var user = await authUsecase.RegisterAsync(req.Email, req.Password);
                var response = new AuthResponse(
                    "Пользователь успешно зарегистрирован",
                    new UserResponse(user.Id, user.Email, user.CreatedAt)
                );
                return Results.Created("/api/v1/auth/register", response);
            }
            catch (InvalidOperationException ex) when (ex.Message == "user already exists")
            {
                return Results.Conflict(new { error = ex.Message });
            }
            catch (Exception ex)
            {
                return Results.Json(new { error = ex.Message }, statusCode: 500);
            }
        })
        .WithName("Register")
        .Accepts<RegisterRequest>("application/json");

        // POST /api/v1/auth/login
        // Принимает email + password, проверяет через BCrypt.
        // 200 — вход выполнен; 401 — неверные данные.
        app.MapPost("/api/v1/auth/login", async (RegisterRequest req) =>
        {
            if (string.IsNullOrEmpty(req.Email) || string.IsNullOrEmpty(req.Password))
                return Results.BadRequest(new { error = "email and password required" });

            try
            {
                var user = await authUsecase.LoginAsync(req.Email, req.Password);
                var response = new AuthResponse(
                    "Добро пожаловать!",
                    new UserResponse(user.Id, user.Email, user.CreatedAt)
                );
                return Results.Json(response);
            }
            catch (UnauthorizedAccessException)
            {
                return Results.Json(new { error = "invalid email or password" }, statusCode: 401);
            }
            catch (Exception ex)
            {
                return Results.Json(new { error = ex.Message }, statusCode: 500);
            }
        })
        .WithName("Login")
        .Accepts<RegisterRequest>("application/json");

        // DELETE /api/v1/auth/delete
        // Принимает email + password — проверяет личность и удаляет аккаунт.
        // 200 — удалён; 401 — неверные данные.
        app.MapDelete("/api/v1/auth/delete", async (RegisterRequest req) =>
        {
            if (string.IsNullOrEmpty(req.Email) || string.IsNullOrEmpty(req.Password))
                return Results.BadRequest(new { error = "email and password required" });

            try
            {
                await authUsecase.DeleteByCredentialsAsync(req.Email, req.Password);
                return Results.Json(new { message = "Пользователь успешно удалён" });
            }
            catch (UnauthorizedAccessException)
            {
                return Results.Json(new { error = "invalid email or password" }, statusCode: 401);
            }
            catch (Exception ex)
            {
                return Results.Json(new { error = ex.Message }, statusCode: 500);
            }
        })
        .WithName("DeleteUser")
        .Accepts<RegisterRequest>("application/json");
    }
}

// RegisterRequest — тело запросов register / login / delete.
// Email и Password приходят в JSON-теле (не в заголовке).
public record RegisterRequest(
    [property: JsonPropertyName("email")]    string Email,
    [property: JsonPropertyName("password")] string Password
);
