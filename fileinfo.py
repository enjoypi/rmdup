"""This is the example module.

This module does stuff.
"""

from hashlib import sha256
from pathlib import Path
from zlib import adler32


class FileInfo(object):
    """
    docstring
    """
    pass

    def __new__(cls, path: Path, filename: str):
        """
        docstring
        """
        filepath = path.joinpath(filename).resolve()
        if not filepath.is_file():
            return None

        stat = filepath.stat()
        if stat.st_size <= 0:
            return UserWarning('Empty file: ' + str(filepath))

        checksum = None
        try:
            with filepath.open(mode='rb') as f:
                checksum = adler32(f.read(8192))
        except OSError as e:
            return e

        instance = super().__new__(cls)
        instance._hash = checksum
        instance.path = filepath
        instance.stat = stat
        return instance

    def __init__(self, _path: Path, _filename: str):
        """
        docstring
        """
        self._digest = None
        self.uri = str(self.path)
        pass

    def __eq__(self, other):
        return self.stat.st_size == other.stat.st_size \
            and self.uri != other.uri and hash(self) == hash(other) \
            and self.digest() == other.digest()

    def __hash__(self):
        return self._hash

    def __repr__(self):
        return self.uri

    def digest(self):
        """
        docstring
        """
        if self._digest is None:
            self._digest = sha256(self.path.read_bytes()).hexdigest()
        return self._digest
