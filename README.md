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

commands:
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

events:
  - name: UserCreated
    type: user.created
    state: true
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

reactors:
  - name: WelcomeEmailReactor
    description: Sends a welcome email to the user when created
    type: local
    reactsTo: user.created

projections:
  - name: UserProjection
    description: Creates a user projection
    type: local
    reactsTo:
      - user.created
      - user.updated
```

