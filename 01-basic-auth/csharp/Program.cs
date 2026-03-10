using BasicAuth.Delivery;
using BasicAuth.Domain;
using BasicAuth.Repository;
using BasicAuth.UseCase;
using BasicAuth.Delivery.Swagger;

// Entry point — C# / ASP.NET Core Minimal API.
//
// Здесь:
// - регистрируем зависимости (Repository → UseCase)
// - подключаем Swagger для тестирования
// - маппим endpoints через MapAuthEndpoints
// - запускаем сервер на :8080
//
// BasicAuth middleware удалён: авторизация теперь выполняется
// непосредственно в /login и /delete через тело запроса.
var builder = WebApplication.CreateBuilder(args);

builder.Services.AddSingleton<IUserRepository, MemoryUserRepository>();
builder.Services.AddSingleton<AuthUsecase>();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen(c =>
{
    c.SchemaFilter<ExamplesSchemaFilter>();
});

var app = builder.Build();

var authUsecase = app.Services.GetRequiredService<AuthUsecase>();

// Swagger включаем всегда — удобно для тестирования в учебных проектах.
app.UseSwagger();
app.UseSwaggerUI(c =>
{
    // Swagger UI открывается по /swagger
    c.SwaggerEndpoint("/swagger/v1/swagger.json", "Basic Auth API v1");
    c.RoutePrefix = "swagger";
});

app.MapAuthEndpoints(authUsecase);

Console.WriteLine("Basic Auth server (C# .NET) running on http://localhost:8080");
Console.WriteLine("Swagger UI: http://localhost:8080/swagger");
app.Run("http://localhost:8080");
