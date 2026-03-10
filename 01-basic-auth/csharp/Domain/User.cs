namespace BasicAuth.Domain;

public class User
{
    public required string Id { get; set; }
    public required string Email { get; set; }
    public required string Password { get; set; }
    public DateTime CreatedAt { get; set; }
}
