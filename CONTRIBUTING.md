# Contributing to oscal-sdk-go

Thank you for your interest in `oscal-sdk-go`!
Our project welcomes external contributions.

## How to Contribute

To contribute code or documentation, please submit a [pull request](https://github.com/oscal-compass/oscal-sdk-go/pulls).

It is preferable that a pull request relates to an existing issue. If you find a bug or want to suggest a feature, please submit
a GitHub [issue](https://github.com/oscal-compass/oscal-sdk-go/issues) first.
This is not required for minor changes like typos in documentation.

## Development

### Run tests

```bash
make test-unit
```

### Format and Style

**Requires [`golangci-lint`](https://golangci-lint.run/welcome/quick-start/)**

```bash
make format
# For issue identification
make vet
# Linting
make lint
```

## Legal

Each source file must include a license header for the Apache
Software License 2.0. Using the SPDX format is the simplest approach.

e.g.
```text
/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/
```

### Sign your work

We have tried to make it as easy as possible to make contributions. This
applies to how we handle the legal aspects of contribution.

We use the same approach - the [Developer's Certificate of Origin 1.1 (DCO)](https://oscal-compass.github.io/compliance-trestle/latest/contributing/DCO/) - that the Linux® Kernel [community](https://developercertificate.org/)
uses to manage code contributions.

We simply ask that when submitting a patch for review, the developer
must include a sign-off statement in the commit message.

Here is an example Signed-off-by line, which indicates that the
submitter accepts the DCO:

```text
Signed-off-by: John Doe <john.doe@example.com>
```

You can include this automatically when you commit a change to your
local git repository using the following command:

```bash
git commit --signoff
```

Note that DCO signoff is enforced by [DCO bot](https://github.com/probot/dco). Missing DCO's will be required to be rebased
with a signed off commit before being accepted.
