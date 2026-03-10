using BasicAuth.Delivery;
using BasicAuth.Domain;
using BasicAuth.Repository;
using BasicAuth.UseCase;
using BasicAuth.Delivery.Middleware;

// Entry point (C# / ASP.NET Core).
//
// Здесь мы:
// - регистрируем зависимости (Repository → UseCase)
// - подключаем Basic Auth middleware для защищённых путей
// - маппим endpoints и запускаем сервер на 8080
var builder = WebApplication.CreateBuilder(args);

builder.Services.AddSingleton<IUserRepository, MemoryUserRepository>();
builder.Services.AddSingleton<AuthUsecase>();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

var app = builder.Build();

var authUsecase = app.Services.GetRequiredService<AuthUsecase>();
app.UseMiddleware<BasicAuthMiddleware>(authUsecase);

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.MapAuthEndpoints(authUsecase);

Console.WriteLine("Basic Auth server (C# .NET) running on http://localhost:8080");
app.Run("http://localhost:8080");
