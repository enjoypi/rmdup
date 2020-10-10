import argparse
import os

from fileinfo import FileInfo, valid
from pathlib import Path


def tidy(directory, collection):
    for root, _, files in os.walk(directory):
        path = Path(root)
        for filename in files:
            if not valid(path, filename):
                continue
            fi = FileInfo(path, filename)
            collection.setdefault(fi, []).append(str(fi))


def main():
    parser = argparse.ArgumentParser(
        description='remove duplicate files', prog='rmdup')
    parser.add_argument('directories', metavar='directory', type=str, nargs='+',
                        help='the directory to find duplicate')

    args = parser.parse_args()
    files = {}
    for directory in args.directories:
        tidy(directory, files)

    for value in files.values():
        if len(value) > 1:
            print(value)


if __name__ == "__main__":
    main()
