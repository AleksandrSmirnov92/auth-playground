using BasicAuth.Delivery.Contracts;
using Microsoft.OpenApi.Any;
using Microsoft.OpenApi.Models;
using Swashbuckle.AspNetCore.SwaggerGen;

namespace BasicAuth.Delivery.Swagger;

// ExamplesSchemaFilter добавляет "человеческие" примеры вместо "string" в Swagger UI.
public class ExamplesSchemaFilter : ISchemaFilter
{
    public void Apply(OpenApiSchema schema, SchemaFilterContext context)
    {
        if (context.Type == typeof(RegisterRequest))
        {
            // Иногда schema-level example не отображается в UI,
            // поэтому задаём примеры на уровне свойств.
            if (schema.Properties.TryGetValue("email", out var email))
                email.Example = new OpenApiString("ivanov@example.com");
            if (schema.Properties.TryGetValue("password", out var password))
                password.Example = new OpenApiString("1234");
            return;
        }

        if (context.Type == typeof(UserResponse))
        {
            if (schema.Properties.TryGetValue("id", out var id))
                id.Example = new OpenApiString("<uuid>");
            if (schema.Properties.TryGetValue("email", out var uEmail))
                uEmail.Example = new OpenApiString("ivanov@example.com");
            if (schema.Properties.TryGetValue("created_at", out var createdAt))
                createdAt.Example = new OpenApiString("2026-03-10T12:00:00Z");
        }
    }
}

