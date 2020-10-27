from fnmatch import fnmatch
import io


class Wildcard:
    def __init__(self, file: str) -> None:
        self._match_list = []
        with io.open(file, mode='rt', encoding='utf-8') as ignore_file:
            for line in ignore_file:
                line = line.strip()
                if len(line) <= 0:
                    continue
                if line[-1] == '/':
                    self._match_list.append(line + '*')
                else:
                    self._match_list.append(line)
            pass

        return

    def match(self, full_path: str) -> bool:
        """
        docstring
        """
        for pattern in self._match_list:
            if fnmatch(full_path, pat=pattern):
                return True
        return False
