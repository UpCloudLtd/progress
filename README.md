# Progress

[![Test](https://github.com/kangasta/progress/actions/workflows/test.yml/badge.svg)](https://github.com/kangasta/progress/actions/workflows/test.yml)

Go module for printing progress/task status of command line applications.

## Installation

To install this library, use `go get`:

```sh
go get github.com/kangasta/progress
```

## Usage

To log progress messages with Progress, you need to

- [Initialize a Progress instance](#initialize-progress)
- [Push messages to the Progress instance](#push-messages)

### Initialize Progress

To initialize Progress, run `progress.NewProgress(...)`. This initializes the internal message store and message renderer.

To start logging the progress messages, call `Start()` to launch goroutine responsible for updating the progress log and handling incoming progress updates.

When done with logging progress updates, call `Stop()` to render the final progress log state and to terminate the goroutine started by `Start()`.

```go
cfg := progress.DefaultOutputConfig
taskLog := progress.NewProgress(cfg)

taskLog.Start()
defer taskLog.Stop()
```

### Push messages

To push messages to the progress log, call `Push(...)`. For example:

```go
taskLog.Push(messages.Update{
    Key:     "error-example",
    Message: "Pushing message to the progress log",
    Status:  messages.MessageStatusStarted,
})
```

An update can contain four different fileds: `Key`, `Message`, `Status`, and `Details`. When updating progress message that does not yet exist, `Message` and `Status` are required.

Field   | Description
------- | -----------
`Key`     | Used to identify the message when pushing further updates. If not given when pushing first update to the message, value from `Message` field is used as `Key`.
`Message` | Text to be outputted in the progress log for the related message.
`Status`  | Status of the message, e.g. `success`, `error`, `warning`. Used to determine status indicator and color. Finished statuses (`success`, `error`, `warning`, `skipped`, `unknown`) are outputted to persistent log and can not be edited anymore.
`Details` | Details to be outputted under finished progress log row, e.g. error message.

Progress messages can be updated while they are in `pending` or `started` states. Note that `pending` messages are not outputted at the moment.

When updating existing progress message, unchanged fields can be omitted. For example:

```go
taskLog.Push(messages.Update{
    Key:     "error-example",
    Status:  messages.MessageStatusError,
    Details: "Error: Message details can be used, for example, to communicate error messages to the user.",
})
```

## Development

To lint the code, run `golangci-lint run`. See its documentation for  [local installation instructions](https://golangci-lint.run/usage/install/#local-installation).

```sh
golangci-lint run
```

To test the code, run `go test ./...`.

```sh
go test ./...
```

To build and run the example application, run:

```sh
go build -o example.bin ./example/example.go
./example.bin
```
