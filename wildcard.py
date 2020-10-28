from fnmatch import fnmatch
from pathlib import Path, PurePath


class Wildcard:
    def __init__(self, dirpath: str = None, file: str = None) -> None:
        self._patterns = []
        if dirpath and file:
            self.append(dirpath, file)

    def append(self, dirpath: str, file: str) -> None:
        try:
            self.parse_file(dirpath, file)
        except FileNotFoundError:
            pass

    def match(self, full_path: str) -> bool:
        """
        docstring
        """
        for pat in self._patterns:
            if full_path == pat:
                return True
            if fnmatch(full_path, pat=pat):
                return True
        return False

    def parse_file(self, dirpath: str, file: str) -> None:
        pats = set(self._patterns)
        pure = PurePath(dirpath)
        with Path(dirpath).joinpath(file).open(
                mode='rt', encoding='utf-8') as ignore_file:
            for line in ignore_file:
                wild_str = line.strip()
                if len(wild_str) <= 0:
                    continue
                pats.add(str(pure.joinpath(wild_str)))
                if wild_str[-1] == '/':
                    pats.add(str(pure.joinpath(wild_str, '*')))
            pass
        self._patterns = list(pats)
        self._patterns.sort()
        return None
