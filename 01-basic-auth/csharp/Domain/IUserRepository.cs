namespace BasicAuth.Domain;

// IUserRepository — контракт хранилища пользователей.
//
// UseCase зависит от интерфейса, поэтому реализацию можно менять
// (in-memory → база данных) без изменения бизнес-логики.
public interface IUserRepository
{
    Task CreateAsync(User user);
    Task<User?> GetByIdAsync(string id);
    Task<User?> GetByEmailAsync(string email);
    Task DeleteAsync(string id);
}
