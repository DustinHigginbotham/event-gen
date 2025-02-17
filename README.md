# event-gen

event-gen is a tool that generates event-sourced code for your application.

## Installation

Use with the `go tool` command (requires go 1.24+):

```bash
go get -tool github.com/DustinHigginbotham/event-gen@latest
```
And then run it:
```bash
go tool event-gen
```

## Configuration

Include a folder in your project root called `event-gen`. Each file in this folder will represent a domain.

Example:

`event-gen/users.yaml`

```yaml

# Creates a User entity
entity:
  name: User
  description: A user
  fields:
    - name: id
      type: string
    - name: first_name
      type: string
    - name: last_name
      type: string
    - name: email
      type: string

# Sets up a list of commands for the entity
commands:

  # A command to create a new user
  - name: CreateUser
    description: Create a new user
    emits: UserCreated
    handler: domains/user/command_create:user
    fields:
      - name: id
        type: string
      - name: first_name
        type: string
      - name: last_name
        type: string
      - name: email
        type: string

  # A command to update a user
  - name: UpdateUser
    description: Updates the name of the user.
    emits: UserUpdated
    handler: domains/user/command_update:user
    fields:
      - name: id
        type: string
      - name: first_name
        type: string
      - name: last_name
        type: string

# Sets up a list of events for the entity (spawned from the commands)
events:

  # Fired from the CreateUser command
  - name: UserCreated
    type: user.created
    state: true # Changes the state of the entity
    handler: domains/user/handlers:user
    description: A user has been created
    fields:
      - name: id
        type: string
      - name: first_name
        type: string
      - name: last_name
        type: string
      - name: email
        type: string

  # Fired from the UpdateUser command
  - name: UserUpdated
    type: user.updated
    handler: domains/user/handlers:user
    state: true
    description: A user has been updated
    fields:
      - name: id
        type: string
      - name: first_name
        type: string
      - name: last_name
        type: string

# Sets up a list of reactors for the entity
reactors:

  # A reactor to send a welcome email to the user when created
  - name: WelcomeEmailReactor
    description: Sends a welcome email to the user when created
    type: local
    reactsTo: user.created

# Sets up a list of projections for the entity
# This is used to keep a user projection up to date
projections:

  # A projection to create a user projection
  - name: UserProjection
    description: Creates a user projection
    type: local
    reactsTo:
      - user.created
      - user.updated
```

