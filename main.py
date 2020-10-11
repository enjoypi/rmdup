"""This is the example module.

This module does stuff.
"""
import argparse
import logging
import os
from fileinfo import FileInfo
from pathlib import Path


def tidy(directory, collection):
    for root, _, files in os.walk(directory):
        path = Path(root)
        for filename in files:
            fi = FileInfo(path, filename)
            if isinstance(fi, FileInfo):
                collection.setdefault(fi, []).append(str(fi))
            else:
                logging.debug(fi)


def main():
    parser = argparse.ArgumentParser(description='remove duplicate files',
                                     prog='rmdup')
    parser.add_argument('-debug',
                        action='store_true',
                        default=False,
                        help='debug')
    parser.add_argument('directories',
                        metavar='directory',
                        type=str,
                        nargs='+',
                        help='the directory to find duplicate')

    args = parser.parse_args()
    # debug = parser.parse_args('-debug')
    # if debug:
    #     logging.basicConfig(level=logging.DEBUG)
    files = {}
    for directory in args.directories:
        tidy(directory, files)

    # for value in files.values():
    #     if len(value) > 1:
    #         print(value)


if __name__ == '__main__':
    main()
