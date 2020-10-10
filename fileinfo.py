from pathlib import PurePath

import hashlib
import zlib


def valid(path, filename):
    fullpath = path.joinpath(filename)
    stat = fullpath.stat()
    if fullpath.is_file() and stat.st_size > 0:
        return True
    return False


class FileInfo(object):
    """
    docstring
    """
    pass

    def __init__(self, path, filename):
        """
        docstring
        """
        fullpath = path.joinpath(filename)

        self._digest = ''
        self._hash = 0
        self.path = fullpath.resolve()
        self.stat = fullpath.stat()
        self.uri = str(self.path)
        pass

    def __eq__(self, other):
        return self.stat.st_size == other.stat.st_size and self.uri != other.uri and hash(self) > 0 and hash(self) == hash(other) and self.digest() == other.digest()

    def __hash__(self):
        if self._hash <= 0:
            try:
                with self.path.open(mode='rb') as f:
                    self._hash = zlib.adler32(f.read(8192))
            except:
                pass
        return self._hash

    def __repr__(self):
        return self.uri

    def digest(self):
        """
        docstring
        """
        if self._digest == '':
            self._digest = hashlib.sha256(self.path.read_bytes()).hexdigest()
        return self._digest
