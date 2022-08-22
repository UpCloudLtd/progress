# Progress

[![Test](https://github.com/UpCloudLtd/progress/actions/workflows/test.yml/badge.svg)](https://github.com/UpCloudLtd/progress/actions/workflows/test.yml)

Go module for printing progress/task status of command line applications.

## Installation

To install this library, use `go get`:

```sh
go get github.com/UpCloudLtd/progress
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
cfg := progress.GetDefaultOutputConfig()
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
`ProgressMessage` | Progress indicator text to be appended into `Message` in TTY terminals, e.g. `128 / 384 kB` or `24 %`. Updating this field will not trigger message write in non-TTY terminals.
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

Use [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) when committing your changes.

To lint the code, run `golangci-lint run`. See its documentation for  [local installation instructions](https://golangci-lint.run/usage/install/#local-installation).

```sh
golangci-lint run
```

To test the code, run `go test ./...`.

```sh
go test ./...
```

To update snapshots, run tests with `UPDATE_SNAPSHOTS` environment variable set to `true`.

```sh
UPDATE_SNAPSHOTS=true go test ./...
```

To build and run the example application, run:

```sh
go build -o example.bin ./example/example.go
./example.bin
```

## Releasing

When releasing a new version:

1. Merge all changes to be included to the `main` branch.
1. Prepare [CHANGELOG.md](./CHANGELOG.md) for the new release:
    1. Add new heading with the correct version (e.g., `## [v2.3.5]`).
    1. Update links at the bottom of the page.
    1. Leave `Unreleased` section at the top empty.
1. Draft a new release in [GitHub releases](https://github.com/UpCloudLtd/progress/releases):
    1. Set the release to create new tag (e.g., `v2.3.5`).
    1. Select the stable branch.
    1. Title the release with the version number (e.g., `v2.3.5`).
    1. In the description of the release, paste the changes from [CHANGELOG.md](./CHANGELOG.md) for this version release.
1. Publish the release when ready.
