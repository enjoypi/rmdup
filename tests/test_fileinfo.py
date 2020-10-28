from tempfile import NamedTemporaryFile
from zlib import adler32

from fileinfo import FileInfo


def test_little_file(monkeypatch):
    data = bytes('''
.git/
*/.git/
*.jpg
*.py
a/*/c/
*/mame0184-64bit/

''',
                 encoding='utf8')
    tmp = NamedTemporaryFile()
    tmp.write(data)
    tmp.flush()
    assert hash(FileInfo(tmp.name)) == adler32(data)
    tmp.close()
