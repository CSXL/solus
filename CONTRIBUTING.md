# Welcome to Solus!

Thank you for joining us on this journey to automating the development process.

## Table of Contents

- [Welcome to Solus!](#welcome-to-solus)
  - [Table of Contents](#table-of-contents)
  - [Executive Summary](#executive-summary)
    - [Mission](#mission)
    - [AI Principles](#ai-principles)
  - [Building, Running, Testing, and Debugging](#building-running-testing-and-debugging)
    - [Environment Secrets](#environment-secrets)
    - [TUI Configuration](#tui-configuration)
    - [Building and Running the Project](#building-and-running-the-project)
    - [Running Tests](#running-tests)
    - [Debugging](#debugging)
  - [Contributing](#contributing)
    - [Code of Conduct](#code-of-conduct)
    - [Contributing Guidelines](#contributing-guidelines)
    - [Types of Contributions](#types-of-contributions)
    - [How to Contribute](#how-to-contribute)
  - [Thank You!](#thank-you)

## Executive Summary

### Mission

To automate the process of creating a project from requirements.

### AI Principles

These are guiding principles specific to the Solus project.
For CSX Labs' core values please refer to our [Business Plan](https://docs.google.com/document/d/1PhPFI1YXRd-XHMvfvRZhFwnqzzdXLTcpo0Kmbw803-I/edit?usp=sharing).

- **Build for Longevity:** AI is a rapidly evolving field that is growing increasingly integrated into our lives. We want to build a platform that will be able to adapt to the changing landscape of AI and will be able to continue to provide value to developers for years to come. At CSX Labs, we make all of our decesions with the future in mind, and this project is no exception.
- **Be Open, Transparent, and Inclusive:** Language models have the capacity to automate our work in tremendous ways. Being the creator of the automating tools is a huge responsibility. By giving _everyone_ the keys to this technology, we mitigate the potential for abuse of power and ensure that the project has a positive impact on our community.
- **Focus on the End-User:** At the end of the day, we are building this project for developers. We will ensure that the project is easy to use and that it ultimatley enhances the field of software development. With AI, there are often concerns of job loss. This tool isn't meant to replace us entirely, but allow us to focus on the more creative aspects of our work.

## Building, Running, Testing, and Debugging

### Environment Secrets

Create an [OpenAI API Key](https://platform.openai.com/account/api-keys) and set it as an environment variable named `OPENAI_API_KEY`.

Create a [Google Cloud Platform](https://cloud.google.com/) account and make a Google Programmable Search Engine API key. You can follow Google's [Custom Search API introduction guide](https://developers.google.com/custom-search/v1/introduction) for for information on obtaining an API key. Set the API key as an environment variable named `GOOGLE_API_KEY`. You will also need to set the `GOOGLE_PROGRAMMABLE_SEARCH_ENGINE_ID` environment variable to the ID of your Programmable Search Engine.

You can also set the environment variables in a `.env` file in the root of the project (see [.env.example](.env.example) for an example .env configuration). The `.env` file is ignored by git, so you can safely store your API key in it.

We are not responsible for any charges incurred by your OpenAI or Google Cloud accounts.

### TUI Configuration

The TUI configuration file ([tui_config.yaml](tui_config.yaml)) contains configuration around the TUI's appearance and behavior. Here is an example configuration file with explanations for each field:

```yaml
debug: bool # Whether or not to enable debug mode. If true, the TUI will print system messages to the terminal.
load_messages_from_file: bool # Whether or not to load messages from a file. If true, the TUI will populate the message history with messages from the file specified in the `saved_messages_file` field.
saved_messages_file: string # The path to the file containing saved messages. This field is only used if `load_messages_from_file` is true.
discovery_message: string # The system message that is sent to the language model to establish the guidelines and context for the conversation.
```

### Building and Running the Project

To run the project, you will need to have [Go](https://go.dev/) and [Make](https://www.gnu.org/software/make/) installed.

To run the project, run `make run` or just `make`, this project uses the terminal AltScreen, so ensure your terminal window is large enough to display the entire TUI.

To build a project binary, run `make build`, and a binary will be written to `solus.out`. You can then run the binary with `./solus.out`.

### Running Tests

Run `make test` to run all the unit tests in the project. We are using VSCode's [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go) for testing, so you can also run tests from within VSCode.

### Debugging

We are using [Delve](https://github.com/go-delve/delve) for debugging integrated with VSCode's [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go). You can run the project in debug mode by pressing `F5` in VSCode.

Word of advice: If you are debugging the TUI manually with Delve you need to run the devserver with `make run` and then attach to the process with `delve attach <pid>`. Another option is to run a headless delve server with `dlv debug --headless .`, connect to the headless server with `dlv connect <address of server>`.

## Contributing

### Code of Conduct

We have a [Code of Conduct](CODE_OF_CONDUCT.md) that we expect all contributors to follow. Please read it before contributing.

Please report unacceptable behavior to [opensource@csxlabs.org](mailto:opensource@csxlabs.org).

### Contributing Guidelines

Follow our [AI Principles](#ai-principles) when contributing to the project. If you have any questions about our code of conduct, guidelines, or operation, feel free to reach out to us at [opensource@csxlabs.org](mailto:opensource@csxlabs.org).

Linting is done with [trunk](https://trunk.io), there are common IDE plugins for it.

### Types of Contributions

We welcome contributions of all kinds. Here are some examples of the types of contributions we are looking for:

- **Code:** Found a bug, want to fill a TODO, or want to enhance our operation? File an issue with a pull request and our open source team will review your changes. If you need help setting up the project or have any questions about the codebase, feel free to reach out to us at [opensource@csxlabs.org](mailto:opensource@csxlabs.org).
- **Documentation:** Have you found a typo in the documentation? Do you have a suggestion for improving the documentation? We would love to hear from you! Feel free to open an issue and/or pull request with your changes.
- **Design:** If you have any suggestions for improving the project's design, feel free to open an issue and/or pull request with your suggestions.

### How to Contribute

1. [Fork](https://docs.github.com/en/get-started/quickstart/fork-a-repo) the repository.
2. [Open an issue](https://docs.github.com/en/issues/tracking-your-work-with-issues/creating-an-issue) describing the changes you would like to make. If you are not sure if your changes are necessary, feel free to open an issue and ask us.
3. If you want to make changes to the repository, [create a pull request](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request) and reference it in your issue. If you are not sure how to make the changes, feel free to open a pull request with your proposed changes and ask us for feedback.

## Thank You!

Thank you for taking the time to read this document. We hope you enjoy using Solus and that it helps you in your work. If you have any questions, comments, or concerns, feel free to reach out to us at [opensource@csxlabs.org](mailto:opensource@csxlabs.org). If you have any business inquiries about Solus or CSX Labs, fill out our [contact form](https://csxlabs.org/#contact) or email us at [info@csxlabs.org](mailto:info@csxlabs.org).
