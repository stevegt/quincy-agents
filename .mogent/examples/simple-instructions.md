# Instructions

## Workflow

1. Read the task carefully
2. Ask questions if unclear
3. Make minimal changes
4. Test before handoff

## Building

```bash
go build -o /tmp/mogent ./tools/mogent
/tmp/mogent build --dry
```

## Testing

```bash
go test ./tools/mogent/...
```

## Committing

- Use imperative mood: "Add feature" not "Added feature"
- One logical change per commit
- Reference TODO files when applicable
