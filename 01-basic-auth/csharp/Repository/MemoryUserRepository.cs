using BasicAuth.Domain;

namespace BasicAuth.Repository;

// In-memory реализация IUserRepository.
//
// Важно:
// - данные живут только в памяти процесса (после перезапуска исчезнут)
// - GetByEmail использует перебор (O(n)) — для учебного примера достаточно
public class MemoryUserRepository : IUserRepository
{
    private readonly Dictionary<string, User> _users = new();

    public Task CreateAsync(User user)
    {
        _users[user.Id] = user;
        return Task.CompletedTask;
    }

    public Task<User?> GetByIdAsync(string id)
    {
        _users.TryGetValue(id, out var user);
        return Task.FromResult(user);
    }

    public Task<User?> GetByEmailAsync(string email)
    {
        var user = _users.Values.FirstOrDefault(u => u.Email == email);
        return Task.FromResult(user);
    }

    public Task DeleteAsync(string id)
    {
        _users.Remove(id);
        return Task.CompletedTask;
    }
}
