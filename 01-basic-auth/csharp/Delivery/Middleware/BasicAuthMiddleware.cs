using BasicAuth.UseCase;
using System.Text;

namespace BasicAuth.Delivery.Middleware;

// BasicAuthMiddleware проверяет HTTP Basic Auth для защищённых эндпоинтов.
//
// Что делает:
// - пропускает публичные пути (например /health и /api/v1/auth/register)
// - для /api/v1/auth/me читает Authorization: Basic ...
// - декодирует Base64(email:password) и вызывает AuthUsecase.LoginAsync
// - при успехе кладёт user.Id в HttpContext.Items["UserId"], чтобы endpoint мог взять его
public class BasicAuthMiddleware
{
    private readonly RequestDelegate _next;
    private readonly AuthUsecase _authUsecase;
    private static readonly PathString[] ProtectedPaths = { "/api/v1/auth/me" };

    public BasicAuthMiddleware(RequestDelegate next, AuthUsecase authUsecase)
    {
        _next = next;
        _authUsecase = authUsecase;
    }

    public async Task InvokeAsync(HttpContext context)
    {
        var path = context.Request.Path;
        if (!ProtectedPaths.Any(p => path.StartsWithSegments(p)))
        {
            await _next(context);
            return;
        }

        var authHeader = context.Request.Headers.Authorization.FirstOrDefault();
        if (string.IsNullOrEmpty(authHeader) || !authHeader.StartsWith("Basic "))
        {
            context.Response.Headers.WWWAuthenticate = "Basic realm=\"Restricted\"";
            context.Response.StatusCode = 401;
            await context.Response.WriteAsJsonAsync(new { error = "Authorization required" });
            return;
        }

        var base64 = authHeader[6..];
        var decoded = Encoding.UTF8.GetString(Convert.FromBase64String(base64));
        var parts = decoded.Split(':', 2);
        if (parts.Length != 2)
        {
            context.Response.Headers.WWWAuthenticate = "Basic realm=\"Restricted\"";
            context.Response.StatusCode = 401;
            await context.Response.WriteAsJsonAsync(new { error = "Invalid credentials" });
            return;
        }

        try
        {
            var user = await _authUsecase.LoginAsync(parts[0], parts[1]);
            context.Items["UserId"] = user.Id;
            await _next(context);
        }
        catch (UnauthorizedAccessException)
        {
            context.Response.Headers.WWWAuthenticate = "Basic realm=\"Restricted\"";
            context.Response.StatusCode = 401;
            await context.Response.WriteAsJsonAsync(new { error = "Invalid credentials" });
        }
    }
}
