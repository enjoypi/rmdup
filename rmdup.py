#!/usr/bin/env python3
"""This is the example module.

This module does stuff.
"""
import argparse
import logging
import os
from collections import OrderedDict
from fileinfo import FileInfo
from pathlib import PurePath

from wildcard import Wildcard


def gen_bat_cmd(filename: str):
    source = PurePath(filename)
    parts = list(source.parts)
    parts.insert(1, 'tmp')
    target = PurePath(*parts)
    return 'ROBOCOPY /MOV "{}" "{}" "{}"\n'.format(str(source.parent),
                                                   str(target.parent),
                                                   source.name)


def gen_sh_cmd(filename: str):
    source = PurePath(filename)
    return 'mkdir -p "/tmp{dir}" && mv "{src}" "/tmp{src}"\n'.format(
        dir=source.parent, src=source)


def generate_script(collection, script, gen_cmd):
    for names in collection:
        if len(names) <= 1:
            continue

        for filename in names:
            script.write(gen_cmd(filename))
        script.write('\n')


def tidy(directory: str, collection: dict):
    ignoring_file = '.rmdupignore'
    wildcard = Wildcard()
    for dirpath, _, filenames in os.walk(directory):
        wildcard.append(dirpath, ignoring_file)
        if wildcard.match(dirpath):
            continue

        path = PurePath(dirpath)
        for filename in filenames:
            if filename == ignoring_file or filename == 'Icon\r':
                continue

            filepath = str(path.joinpath(filename))

            # ignore
            if wildcard.match(filepath):
                continue

            f = FileInfo(filepath)
            if isinstance(f, FileInfo):
                logging.debug(f.uri)
                collection.setdefault(f, []).append(str(f))
            else:
                logging.warning(f)


def main():
    parser = argparse.ArgumentParser(description='remove duplicate files',
                                     prog='rmdup')
    # parser.add_argument('-log.level', default='INFO', help='logging level')
    parser.add_argument('directories',
                        metavar='directory',
                        type=str,
                        nargs='+',
                        help='the directory to find duplicate')

    # level = parser.parse_args('-log.level')
    logging.basicConfig(level='INFO')

    args = parser.parse_args()
    files = OrderedDict()
    for directory in args.directories:
        tidy(directory, files)

    if os.name == 'nt':
        with open('rmdup.bat', mode='wt', encoding='gbk') as script:
            script.write('@ECHO OFF\n')
            generate_script(files.values(), script, gen_bat_cmd)
    else:
        with open('rmdup.sh', mode='wt', encoding='utf-8') as script:
            script.write('#!/bin/sh\n')
            generate_script(files.values(), script, gen_sh_cmd)


if __name__ == '__main__':
    main()
