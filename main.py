from solus import agent

RefactorAgent = agent.RefactorAgent
DeliberationMessage = agent.DeliberationMessage
ExecutorCommand = agent.ExecutorCommand
from colorama import Back as bg
from colorama import Fore as colors
from colorama import Style as styles

colors.RESET = styles.RESET_ALL

# Below is some very inefficient code to refactor
code = '''
# Adds two numbers and reverses it within the same function
def addAndReverse(a, b, alist):
    result = a + b
    i = 0
    j = len(alist) - 1
    # Outdated comment
    while i < j:
        alist[i], alist[j] = alist[j], alist[i]
        i += 1
        j -= 1
    return result, alist
'''


def adjust_lines(code: str) -> str:
  lines = code.split('\n')
  max_length = max(len(line) for line in lines)
  return '\n'.join(line.ljust(max_length) for line in lines)


agent = RefactorAgent(code)
initial_view = agent.resource.view()
agent.refactor()
final_view = agent.resource.view()
steps = agent.context.ephemeral_context


def format_steps(context: list) -> str:
  steps = ''
  for item in context:
    if type(item) == DeliberationMessage:
      steps += f'{colors.BLUE}{styles.BRIGHT}DELIBERATOR: {colors.RESET}{colors.BLUE}{item.message}{colors.RESET}'
    elif type(item) == ExecutorCommand:
      steps += f'{colors.YELLOW}{styles.BRIGHT}EXCECUTOR: {colors.RESET}{colors.YELLOW}{item.command_type} {item.truncate_params(item.params)}{colors.RESET}'
    steps += '\n'
  return steps


print(f'''
{styles.BRIGHT}SUMMARY
_______

Initial Code
{colors.BLACK}{bg.WHITE}{adjust_lines(initial_view)}
{colors.RESET}
{styles.BRIGHT}Refactored Code
{colors.BLACK}{bg.WHITE}{adjust_lines(final_view)}
{colors.RESET}
{styles.BRIGHT}STEPS
_____
{format_steps(steps)}
{colors.RESET}

'''.strip())
