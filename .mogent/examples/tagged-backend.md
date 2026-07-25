# Backend API Style {#team/backend #lang/go}

## REST Patterns

- Use plural nouns: `/users`, `/orders`
- Version in URL: `/v1/users`
- Return proper status codes

## Error Responses

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "User not found"
  }
}
```

## Database

- Use transactions for multi-step ops
- Prefer migrations over manual SQL
- Index foreign keys
