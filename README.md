# Solus
An `AI-assisted` project generator.

See the project [specification](SPECIFICATION.md) for more specific technical details.

## Mission
To automate the process of creating a project from requirements.

## Why?
As a developer, much of our time is spent on interpreting requirements into 
code and then testing it. With Solus, we can reduce the time spent on low-level
tasks and focus on high-level tasks like architecture, design, and testing.

## Requirements
**Vision:** A CLI tool that can generate an end-to-end project from well-defined 
requirements.

### Collecting Requirements
The first step is to collect requirements from the developer. Of course, as 
developers, we don't always have a clear idea of what we want. Solus will 
inquire about the project and then generate a template for the developer to 
fill out.

### Generating a Project
Once the developer has filled out the template, Solus will generate a project
from the requirements. The project will be generated in a directory that the
developer can specify or in the current directory. The AI will be able to query
the internet to get more information about current technologies and APIs, as the
AI has limited knowledge of its training data.

### Test Driven Development
The generated project will follow a test-driven development approach. Solus
will first generate the tests, confirm them with the developer, and then 
generate the code to pass the tests. Writing tests first from a well-defined 
specification helps against the traditional cohesion problems with current 
AI code generation.

### (AI) Code Review
Solus will review the generated code, running linters and other tools to 
ensure that the code is up to standard. The developer will be able to make 
requests to Solus to change the generated code, and Solus will make the 
changes and re-run the tests.

## Task List
This task list is for the Minimum Viable Product (MVP) of Solus. The MVP will 
meet the bare minimum of the requirements listed above and in the specification.
The tasks are being completed according to the [Tracer Bullet Philosophy](https://wiki.c2.com/?TracerBullets)
presented in *[The Pragmatic Programmer](https://dl.acm.org/doi/10.5555/320326)*.
The idea is that we will try to slice through the different parts of the project
to provide a simple set of presentable features that can we can continue to 
build upon.

- [ ] Serialization
  - [x] Implement YAML serialization
  - [ ] Implement JSON serialization
- [ ] Query API
  - [ ] Implement Data Gathering
    - [ ] Implement Colly for scraping
    - [x] Implement Google Custom Search API for general searching
    - [ ] Implement StackExchange API for debugging
    - [x] Implement Wikipedia API for getting topic overviews
  - [ ] Implement AI Processing
    - [ ] Implement gpt-3.5-turbo for extracting relevant information from collected data
    - [ ] Implement gpt-3.5-turbo for serializing data into a standardized format
    - [ ] Implement verification of AI output
- [x] Interface
  - [x] Implement CLI
    - [x] Implement [Bubbletea](https://github.com/charmbracelet/bubbletea) for the TUI
    - [x] Implement [Cobra](https://github.com/spf13/cobra) for command management
  - [x] Implement conversations
    - [x] Implement and wrap the [gpt-3.5-turbo](https://openai.com/blog/introducing-chatgpt-and-whisper-apis) language model
  - [x] Implement configuration with [Viper](https://github.com/spf13/viper)
- [ ] Requirements API
  - [ ] Implement gpt-3.5-turbo for processing requirements from conversations
- [ ] Outline API
  - [ ] Implement gpt-3.5-turbo for generating project outline
    - [ ] Implement gpt-3.5-turbo for revising project outline
  - [ ] Implement file and directory generation
- [ ] Dependency API
  - [ ] Implement Query API for resolving dependencies and their documentation
  - [ ] Implement generation of dependency names and documentation
- [ ] Code API
  - [ ] Implement gpt-3.5-turbo for generating test stubs, code, and comments
    - [ ] Implement context-aware code generation and comments that align with the requirements and dependencies
  - [ ] Implement refactoring of code and incorporation of code standards
- [ ] Debugging API
  - [ ] TODO: Figure out how to implement this.
