# Solus
An `AI-assisted` project generator.

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

## Specification
**Language:** Go

### Query API
A library that requests outside information (from Google, StackOverflow, etc.)
and returns the information in a standardized text format. This API should also
be able to get documentation and specifications, such as an API specification,
library list, or platform documentation.

We will use [Google's Programmable Search Engine](https://programmablesearchengine.google.com/about/), 
in conjunction with their [Custom Search API](https://developers.google.com/custom-search/v1/introduction), 
to find specific information and tech lists. We will use the [StackExchange API](https://api.stackexchange.com/), 
to get information from StackOverflow. We will use [Wikipedia's API](https://www.mediawiki.org/wiki/API:Main_page),
to get overviews of technologies and platforms. To scrape websites, we will use
[Colly](https://go-colly.org/), a batteries-included web scraping framework for
Go.

The role of AI in this component is to take the information requested from the 
APIs and process it into a standardized format for the correct purpose. We will
be using OpenAI's [Codex](https://platform.openai.com/docs/models/codex) to
extract the relevant information from the API responses and return it in JSON.

Request Format:
```json
{
  "query": <query: string>,
  "type": <type: string> - "documentation" | "libraries" | "api-specification" | "overview" | "other"
}
```

Example Requests:
```json
{
  "query": "What is the documentation for the Python requests library?",
  "type": "documentation"
}
```
```json
{
  "query": "What are some libraries in Go for websockets?",
  "type": "libraries"
}
```

Response Format:
```json
{
  "response": <response: string>,
  "type": <type: string> - "documentation" | "libraries" | "api-specification" | "overview" | "other"
}
```

### Interface
A CLI tool that collects requirements from the developer in a conversational 
manner and then generates a project from those requirements. The tool will call
the Query API to help the developer choose the correct technologies and APIs, 
and call the rest of the APIs to generate the project.

### Requirements API
A library that processes a conversation with the developer about requirements and 
generates a YAML file that can be used to generate a project.

Response Format:
```yaml
name: <project name>
mission: <project mission>
requirements:
  - <requirement 1>
    - <requirement 1.1>
    - <requirement 1.2>
    - ...
  - <requirement 2>
  - <requirement 3>
  - ...
```

### Outline API
A library that takes in a YAML file of requirements and generates a YAML project 
outline. The outline will be a list of files and directories that will be 
generated by the project. OpenAI [Codex](https://platform.openai.com/docs/models/codex) 
will be used to generate and revise the outline. 
Each file and directory will have a description that will be generated by
OpenAI [Codex](https://platform.openai.com/docs/models/codex) which contains the
requirements that the file or directory will fulfill.

Response Format:
```yaml
language: <programming language>
objects:
  - directory: <directory name>
    - directory: <directory name>
      - ...
    - file: <file name>
      - requirement: <requirement>
      - requirement: <requirement>
      - ...
    - file: <file name>
      - requirement: <requirement>
      - requirement: <requirement>
      - ...
    - ...
  - directory: <directory name>
    - ...
  - ...
```

### Dependency API
A library that takes in a YAML file of requirements and uses the Query API to
find the dependencies for the project as well as their documentation.

Response Format:
```yaml
language: <programming language>
dependencies:
  - <dependency 1>
    - name: <dependency name>
    - documentation: <dependency documentation>
    - ...
  - <dependency 2>
    - name: <dependency name>
    - documentation: <dependency documentation>
rest-apis:
  - <rest-api 1>
    - name: <rest-api name>
    - documentation: <rest-api documentation>
    - specifcation: <rest-api specification>
    - ...
  - <rest-api 2>
    - name: <rest-api name>
    - documentation: <rest-api documentation>
    - ...
```

### Code API
A library that takes in requirements that a file or directory will fulfill and
generates the code needed to fulfill those requirements. The Code API will 
generate the tests first, then the code to pass the tests. The Code API will
use dependency documentation to generate the code.

The code API can also Refactor the code to make it more readable and less 
repetitive. The code API will write clear _useful_ comments to explain the code.

Process:
* Generate test stubs with documenting comments
* Generate tests
* Generate code stubs with documenting comments
* Generate code to pass tests
* Refactor code
* Repeat

### Debugging API
AIs don't make bug-free code just like humans. The Debugging API will be used to
run tests and linters on the generated code, fixing any bugs that are found with
AI.

## Task List
- [ ] Serialization
  - [ ] Implement YAML serialization
  - [ ] Implement JSON serialization
- [ ] Query API
  - [ ] Implement Data Gathering
    - [ ] Implement Colly for scraping
    - [ ] Implement Google Custom Search API for general searching
    - [ ] Implement StackExchange API for debugging
    - [ ] Implement Wikipedia API for getting topic overviews
  - [ ] Implement AI Processing
    - [ ] Implement Codex for extracting relevant information from collected data
    - [ ] Implement Codex for serializing data into a standardized format
    - [ ] Implement verification of AI output
- [ ] Interface
  - [ ] Implement CLI
    - [ ] Implement [Bubbletea](https://github.com/charmbracelet/bubbletea) for the TUI
    - [ ] Implement [Cobra](https://github.com/spf13/cobra) for command management
  - [ ] Implement conversations
    - [ ] Implement and wrap the [gpt-3.5-turbo](https://openai.com/blog/introducing-chatgpt-and-whisper-apis) language model
- [ ] Requirements API
  - [ ] Implement Codex for processing requirements from conversations
- [ ] Outline API
  - [ ] Implement Codex for generating project outline
    - [ ] Implement Codex for revising project outline
  - [ ] Implement file and directory generation
- [ ] Dependency API
  - [ ] Implement Query API for resolving dependencies and their documentation
  - [ ] Implement generation of dependency names and documentation
- [ ] Code API
  - [ ] Implement Codex for generating test stubs, code, and comments
    - [ ] Implement context-aware code generation and comments that align with the requirements and dependencies
  - [ ] Implement refactoring of code and incorporation of code standards
- [ ] Debugging API
  - [ ] TODO: Figure out how to implement this.
