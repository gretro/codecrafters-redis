["Build Your Own Redis" Challenge](https://codecrafters.io/challenges/redis).

In this challenge, you'll build a toy Redis clone that's capable of handling
basic commands like `PING`, `SET` and `GET`. Along the way we'll learn about
event loops, the Redis protocol and more.

**Note**: If you're viewing this repo on GitHub, head over to
[codecrafters.io](https://codecrafters.io) to try the challenge.

# Requirements

- Go 1.24+
- [Task](https://taskfile.dev/)

## Available commands

| Command | Description |
| --- | --- |
| `task serve` | Runs the back-end project |
| `task test:unit` | Runs the unit tests |
| `task test:int` | Runs the integration tests. Note the server must be running at the same time. |
| `task lint` | Runs the linter |
| `task lint:fix` | Auto-fixes the lint errors that can be fixed |