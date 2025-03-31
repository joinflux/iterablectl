# iterablectl

<!--toc:start-->

- [Installation](#installation)
- [Usage](#usage)
- [Data File Format](#data-file-format)
- [Available Commands](#available-commands)
- [Global Flags](#global-flags)
- [License](#license)

<!--toc:end-->

A command-line tool for interfacing with the Iterable API.

## Installation

```bash
go install github.com/joinflux/iterablectl@latest
```

## Usage

```bash
# Set your Iterable API key
export ITERABLE_API_KEY="your_api_key_here"

# Get a user with explicit API key
iterablectl users get --api-key=your_api_key_here --email=user@example.com

# Get a user by email
iterablectl users get --email=user@example.com

# Get a user with JSON output
iterablectl users get --email=user@example.com --format=json

# Update a user
iterablectl users update --email=user@example.com --data-field=firstName=John --data-field=lastName=Doe

# Update a user with a JSON file containing data fields
iterablectl users update --email=user@example.com --data-file=user_data.json
```

## Data File Format

When updating a user with `--data-file`, the JSON file should contain a flat object of data fields:

```json
{
  "firstName": "John",
  "lastName": "Doe",
  "custom": {
    "favoriteColor": "blue"
  }
}
```

## Available Commands

- `users update` - Update a user in Iterable
- `users get` - Get a user from Iterable by email

## Global Flags

- `--api-key, -k` - Iterable API key (required unless ITERABLE_API_KEY environment variable is set)

## License

MIT
