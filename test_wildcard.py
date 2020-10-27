import unittest
from unittest.mock import mock_open, patch

from wildcard import Wildcard


class IgnoringTestCase(unittest.TestCase):

    # Only use setUp() and tearDown() if necessary

    def setUp(self):
        read_data = '''
.git/
*/.git/
*.jpg
*.py
a/*/c/
*/mame0184-64bit/

'''
        with patch('io.open', mock_open(read_data=read_data)):
            self.wildcard = Wildcard('.', '.rmdupignore')

    def tearDown(self):
        pass

    def test_contain(self):
        assert self.wildcard.match('.git/')
        assert self.wildcard.match('.git/a/b/c')
        assert self.wildcard.match('/a/b/c/.git/')
        assert self.wildcard.match('/a/b/c/.git/d/e/f')
        assert self.wildcard.match('/a/b/c/d.e.f.jpg')
        # Test feature one.
        pass


if __name__ == '__main__':
    unittest.main()
