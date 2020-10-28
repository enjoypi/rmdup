from unittest.mock import mock_open, patch

from wildcard import Wildcard


def test_wildcard(monkeypatch):
    read_data = '''
.git/
*/.git/
*.jpg
*.py
a/*/c/
*/mame0184-64bit/

'''
    with patch('io.open', mock_open(read_data=read_data)):
        wildcard = Wildcard('.', '.rmdupignore')
        assert wildcard.match('.git')
        assert wildcard.match('.git/a/b/c')
        assert wildcard.match('/a/b/c/.git/')
        assert wildcard.match('/a/b/c/.git/d/e/f')
        assert wildcard.match('/a/b/c/d.e.f.jpg')
