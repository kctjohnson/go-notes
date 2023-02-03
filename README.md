# GoNotes Server

## Installation

`go mod tidy`
`go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
`make migrate`
`go build ./cmd/gonotes-server/.`
Move the server to any bin folder, then create a systemctl file for it so it always runs.
Make sure that you've got the port forwarded, otherwise the client won't be able to connect remotely.
