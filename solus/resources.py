class Text():

  def __init__(self, text: str = '') -> None:
    self.text = text

  def get_text(self) -> str:
    return self.text

  def text_by_line(self) -> list[str]:
    '''
    Splits text by line number.
    '''
    return self.get_text().split('\n')

  def join_by_line(self, text_lines: list[str]) -> str:
    self.text = '\n'.join(text_lines)
    return self.get_text()

  def clear(self, replacement: str = '') -> str:
    self.text = replacement
    return self.get_text()

  def insert(self, text: str = '', line_no: int = 1) -> str:
    '''
    Insert text after line number. Line numbers start at 1.
    '''
    lines = self.text_by_line()
    adjusted_line_no = line_no  # We want the insert to be offset by one, so we don't have to change the line number
    lines.insert(adjusted_line_no, text)
    self.join_by_line(lines)
    return self.get_text()

  def delete(self, start_line: int = 1, end_line: int = -1) -> str:
    '''
    Deletes text from start_line to end_line (inclusive). If end_line isn't specified,
    it will just delete the line no. given.
    '''
    lines = self.text_by_line()
    adjusted_start = start_line - 1
    if end_line == -1:
      del lines[adjusted_start]
    else:
      lines_to_delete = [*range(adjusted_start, end_line)]
      for i in sorted(lines_to_delete, reverse=True):
        del lines[i]
    self.join_by_line(lines)
    return self.get_text()

  def edit(self,
           replacement: str = '',
           start_line: int = 1,
           end_line: int = -1) -> str:
    self.delete(start_line, end_line)
    self.insert(replacement, start_line - 1)
    return self.get_text()

  def view(self) -> str:
    '''
    Adds line numbers to every line and returns it in a string. Similar to how you would
    display lines in an IDE.
    '''
    lines = self.text_by_line()
    view = ''
    for no, line in enumerate(lines):
      view += f'{no+1}|{line}\n'
    return view

  def __repr__(self) -> str:
    return self.get_text()
