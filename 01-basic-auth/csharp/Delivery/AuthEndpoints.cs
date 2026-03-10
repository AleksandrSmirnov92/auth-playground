using BasicAuth.UseCase;

namespace BasicAuth.Delivery;

public static class AuthEndpoints
{
    public static void MapAuthEndpoints(this IEndpointRouteBuilder app, AuthUsecase authUsecase)
    {
        app.MapGet("/health", () => Results.Json(new { status = "ok" }));

        app.MapPost("/api/v1/auth/register", async (RegisterRequest req) =>
        {
            if (string.IsNullOrEmpty(req.Email) || string.IsNullOrEmpty(req.Password))
                return Results.BadRequest(new { error = "email and password required" });

            try
            {
                var user = await authUsecase.RegisterAsync(req.Email, req.Password);
                return Results.Created("/api/v1/auth/me", user);
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

        app.MapGet("/api/v1/auth/me", async (HttpContext ctx) =>
        {
            var userId = ctx.Items["UserId"] as string;
            if (string.IsNullOrEmpty(userId))
                return Results.Json(new { error = "unauthorized" }, statusCode: 401);

            var user = await authUsecase.GetUserByIdAsync(userId);
            if (user == null)
                return Results.NotFound(new { error = "user not found" });

            return Results.Json(user);
        })
        .WithName("Me");

        app.MapDelete("/api/v1/auth/me", async (HttpContext ctx) =>
        {
            var userId = ctx.Items["UserId"] as string;
            if (string.IsNullOrEmpty(userId))
                return Results.Json(new { error = "unauthorized" }, statusCode: 401);

            await authUsecase.DeleteUserByIdAsync(userId);
            return Results.NoContent();
        })
        .WithName("DeleteUser");
    }
}

public record RegisterRequest(string Email, string Password);
