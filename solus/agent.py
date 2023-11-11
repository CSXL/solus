import re

from . import model, resources


class AgentContext():
  guiding_principles = '''
Core Values
1. Focus on the user, and all else will follow.
2. Embrace risk to compound success.
3. Integrity and transparency build longevity.
4. Health is wealth.
  '''.strip()
  task_metadata = '''
Role: Refactoring agent.
Goal: Optimize code with idiomatic practices, guiding principles, and performance
constraints.
  '''.strip()
  implementation_standards = '''
Follow Robert C. Martin's clean code practices.
  '''.strip()
  available_commands = '''
INSERT COMMAND
___
INSERT <line> |
<text>
END

Inserts code after the line specified. Example:
Resource View:
1:I work for:
2:The best company
INSERT 1
CSX Labs
END
(New Resource View)
1:I work for:
2:CSX Labs
3:The best company
___

DELETE COMMAND
___
DELETE <start_line> <optional: end_line> END

Deletes the line(s) specified.
___

EDIT COMMAND
___
EDIT <start_line> <optional: end_line> |
<replacement text>
END

Deletes the line(s) specified and then inserts the replacement text at the location
specified. Example:
Resource View:
1:I work for:
2:CSX Labs
3:The best company
EDIT 2 3 |
CSX Labs
The Best Company
In the Universe
END
(New Resource View)
1:I work for:
2:CSX Labs
3:The Best Company
4:In the Universe
___
  '''.strip()

  def __init__(self):
    self.ephemeral_context = []
    self.resource_view = ''

  def get_ephemeral_context(self) -> str:
    context = ''
    for item in self.ephemeral_context:
      context += str(item).strip() + '\n'
    return context

  def append_context(self, item):
    self.ephemeral_context.append(item)

  def update_resource_view(self, view: str):
    self.resource_view = view

  def get_deliberator_context(self) -> str:
    return f'''
You are the deliberator, your role is to reason on the information provided and describe
how to execute the task provided. Your input will then be passed on to an executor that will
take your chain of human reasoning and convert it to commands. Ensure to be short, reasonable, and concise in your explanation.
The executor's responses are inserted into your context, you can give the code to the executor, but don't write "EXECUTOR:".
If you feel the code is sufficient to the task completed, say that the task is completed.
Here are the guiding principles for your organization:
{self.guiding_principles}
Here is the task:
{self.task_metadata}
Here are the implementation standards:
{self.implementation_standards}
Ephemeral Context:
{self.get_ephemeral_context()}
Resource View:
{self.resource_view}
    '''.strip()

  def get_executor_context(self) -> str:
    return f'''
You are the executor, your role is to process the deliberations given by the deliberator and give serializable
commands from the information provided. DO NOT PROVIDE ANY EXPLANATIONS, your output
will be directly processed by a parser. You must only respond with a single command, so choose it wisely.
Here is the task:
{self.task_metadata}
Here are the implementation standards for your code:
{self.implementation_standards}
Here are your available commands:
{self.available_commands}
Ephemeral Context:
{self.get_ephemeral_context()}
Resource View:
{self.resource_view}
    '''.strip()


class DeliberationMessage():

  def __init__(self, message):
    self.message = message

  def __str__(self):
    return f'DELIBERATOR: {self.message}'


class ExecutorCommand():

  def __init__(self, command_type: str, params: dict):
    self.command_type = command_type
    self.params = params

  def truncate_params(self, params: dict) -> dict:
    truncated = {}
    if params == None:
      return ''
    for key, value in params.items():
      if type(value) == str and len(value) > 100:
        value = f'{value[:50]}...{value[-50:]}'
      truncated[key] = value
    return truncated

  def __str__(self):
    return f'EXECUTOR: ```{self.command_type} {self.truncate_params(self.params)}```'


class RefactorAgent():

  def __init__(self, code: str):
    self.resource = resources.Text(code.strip())  # Simulates a file
    self.model = model.LanguageModel()
    self.context = AgentContext()
    self.update_context()

  def update_context(self, item=None, refresh_resource_view=True):
    if refresh_resource_view:
      self.context.update_resource_view(self.resource.view())
    if item:
      self.context.append_context(item)

  def deliberate(self) -> str:
    deliberation = self.model.complete(self.context.get_deliberator_context())
    self.update_context(DeliberationMessage(deliberation))
    return deliberation

  def execute(self) -> dict:
    execution = self.model.complete(self.context.get_executor_context())
    cmd = self.parse_execute_command(execution)
    if cmd['type'] == 'DELETE':
      params = cmd['params']
      self.resource.delete(params['start_line'], params['end_line'])
    elif cmd['type'] == 'INSERT':
      params = cmd['params']
      self.resource.insert(params['text'], params['line'])
    elif cmd['type'] == 'EDIT':
      params = cmd['params']
      self.resource.edit(params['text'], params['start_line'],
                         params['end_line'])

    self.update_context(ExecutorCommand(cmd['type'], cmd['params']))
    return cmd

  def refactor(self, num_iterations=3):
    for i in range(num_iterations):
      try:
        self.deliberate()
        self.execute()
      except:
        pass  # In case the model hits a token limit

  def parse_execute_command(self, command: str) -> dict:
    insert_pattern = r"(INSERT\s[0-9]+\s\|\s.*\sEND)"
    delete_pattern = r"(DELETE\s[0-9]+(\s[0-9]+)?\s?END)"
    edit_pattern = r"(EDIT\s[0-9]+(\s[0-9]+)?\s\|\s.*\sEND)"

    dict_command = {'type': None, 'params': None}

    insert_match = re.search(insert_pattern, command, re.DOTALL)
    delete_match = re.search(delete_pattern, command)
    edit_match = re.search(edit_pattern, command, re.DOTALL)

    if insert_match:
      dict_command['type'] = 'INSERT'
      params = insert_match.group().split("|")
      dict_command['params'] = {
          'line': int(params[0].split()[1]),
          'text': params[1].strip("\nEND")
      }

    elif delete_match:
      dict_command['type'] = 'DELETE'
      params = delete_match.group().split()
      dict_command['params'] = {
          'start_line': int(params[1]),
          'end_line': int(params[2]) if len(params) > 3 else None
      }

    elif edit_match:
      dict_command['type'] = 'EDIT'
      params = edit_match.group().split("|")
      line_params = params[0].split()
      dict_command['params'] = {
          'start_line': int(line_params[1]),
          'end_line': int(line_params[2]) if len(line_params) > 2 else None,
          'text': params[1].strip("\nEND")
      }

    return dict_command
