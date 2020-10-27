from fnmatch import fnmatch
from pathlib import Path
from typing import List


class Wildcard:
    def __init__(self, dirpath: str, file: str) -> None:
        self._patterns = []
        try:
            self._patterns = parse_file(dirpath, file)
        except FileNotFoundError:
            pass

    def match(self, full_path: str) -> bool:
        """
        docstring
        """
        for pat in self._patterns:
            if fnmatch(full_path, pat=pat):
                return True
        return False


def parse_file(dirpath: str, file: str) -> List[str]:
    ret = []
    path = Path(dirpath)
    with path.joinpath(file).open(mode='rt', encoding='utf-8') as ignore_file:
        for line in ignore_file:
            wild_str = line.strip()
            if len(wild_str) <= 0:
                continue
            if wild_str[-1] == '/':
                ret.append(path.joinpath(wild_str))
                ret.append(path.joinpath(wild_str, '*'))
            else:
                ret.append(path.joinpath(wild_str))
        pass

    return ret
