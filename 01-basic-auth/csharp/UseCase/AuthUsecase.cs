using BasicAuth.Domain;

namespace BasicAuth.UseCase;

// AuthUsecase — бизнес-логика авторизации (UseCase слой).
//
// Здесь нет HTTP: только правила приложения:
// - RegisterAsync: email уникален, пароль хешируем (BCrypt), сохраняем пользователя
// - LoginAsync: проверяем email + пароль (BCrypt.Verify)
public class AuthUsecase
{
    private readonly IUserRepository _userRepository;

    public AuthUsecase(IUserRepository userRepository)
    {
        _userRepository = userRepository;
    }

    public async Task<object> RegisterAsync(string email, string password)
    {
        var existing = await _userRepository.GetByEmailAsync(email);
        if (existing != null)
            throw new InvalidOperationException("user already exists");

        var hashedPassword = BCrypt.Net.BCrypt.HashPassword(password);
        var user = new User
        {
            Id = Guid.NewGuid().ToString(),
            Email = email,
            Password = hashedPassword,
            CreatedAt = DateTime.UtcNow
        };

        await _userRepository.CreateAsync(user);
        return new { user.Id, user.Email, user.CreatedAt };
    }

    public async Task<User> LoginAsync(string email, string password)
    {
        var user = await _userRepository.GetByEmailAsync(email);
        if (user == null || !BCrypt.Net.BCrypt.Verify(password, user.Password))
            throw new UnauthorizedAccessException("invalid email or password");
        return user;
    }

    public async Task<object?> GetUserByIdAsync(string id)
    {
        var user = await _userRepository.GetByIdAsync(id);
        return user == null ? null : new { user.Id, user.Email, user.CreatedAt };
    }

    public async Task DeleteUserByIdAsync(string id)
    {
        await _userRepository.DeleteAsync(id);
    }
}
