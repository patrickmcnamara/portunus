# Portunus - a password manager

Portunus is a very basic CLI password manager.

## Installation

Run `go get -u github.com/patrickmcnamara/portunus`.
And then `portunus`, assuming your $GOPATH is in your $PATH.

## Usage

1. Create a portunus vault with `portunus vlt`.
2. Add credentials with `portunus set NAME` or `portunus new NAME`. The former takes the password on the command line, the latter generates a secure password for you.
3. View credentials with `portunus get NAME`.

## Licence

Licenced under the EUPL-1.2.
