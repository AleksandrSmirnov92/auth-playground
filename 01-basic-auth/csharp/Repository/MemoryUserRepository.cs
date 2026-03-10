using BasicAuth.Domain;

namespace BasicAuth.Repository;

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
