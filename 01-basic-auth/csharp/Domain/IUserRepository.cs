namespace BasicAuth.Domain;

public interface IUserRepository
{
    Task CreateAsync(User user);
    Task<User?> GetByIdAsync(string id);
    Task<User?> GetByEmailAsync(string email);
    Task DeleteAsync(string id);
}
