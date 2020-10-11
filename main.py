"""This is the example module.

This module does stuff.
"""
import argparse
import logging
import os
from fileinfo import FileInfo
from pathlib import Path, PurePath


def gen_bat_cmd(filename: str):
    source = PurePath(filename)
    parts = list(source.parts)
    parts.insert(1, 'tmp')
    target = PurePath(*parts)
    return 'ROBOCOPY /MOV "{}" "{}" "{}"\n'.format(str(source.parent),
                                                   str(target.parent),
                                                   source.name)


def generate_bat(collection: [str], script_name: str, gen_cmd):
    bat = open(script_name, mode='wt', encoding='gbk')

    for names in collection:
        if len(names) <= 1:
            continue

        for filename in names:
            bat.write(gen_cmd(filename))
        bat.write('\n')


def tidy(directory: str, collection: dict):
    for dirpath, _, filenames in os.walk(directory):
        path = Path(dirpath)
        for filename in filenames:
            f = FileInfo(path, filename)
            if isinstance(f, FileInfo):
                collection.setdefault(f, []).append(str(f))
            else:
                logging.debug(f)


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

    if os.name == 'nt':
        generate_bat(files.values(), 'rmdup.bat', gen_bat_cmd)


if __name__ == '__main__':
    main()
