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

    public async Task<User> RegisterAsync(string email, string password)
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
        return user;
    }

    public async Task<User> LoginAsync(string email, string password)
    {
        var user = await _userRepository.GetByEmailAsync(email);
        if (user == null || !BCrypt.Net.BCrypt.Verify(password, user.Password))
            throw new UnauthorizedAccessException("invalid email or password");
        return user;
    }

    public async Task<User?> GetUserByIdAsync(string id)
    {
        var user = await _userRepository.GetByIdAsync(id);
        return user;
    }

    public async Task DeleteUserByIdAsync(string id)
    {
        await _userRepository.DeleteAsync(id);
    }

    // DeleteByCredentialsAsync — проверяет email+пароль и удаляет пользователя.
    // Повторно использует LoginAsync, чтобы не дублировать логику проверки пароля.
    public async Task DeleteByCredentialsAsync(string email, string password)
    {
        var user = await LoginAsync(email, password);
        await DeleteUserByIdAsync(user.Id);
    }
}
