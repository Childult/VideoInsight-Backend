import sys, os
import json
import torch


def write_json(splits, save_path):
    if not os.path.exists(os.path.dirname(save_path)):
        os.mkdir(os.path.dirname(save_path))

    with open(save_path, 'w') as f:
        json.dump(splits, f, indent=4, separators=(', ', ': '))


def read_json(fpath):
    with open(fpath, 'r') as f:
        obj = json.load(f)
    return obj


def save_checkpoint(state, fpath='checkpoint.pth.tar'):
    if not os.path.exists(os.path.dirname(fpath)):
        os.mkdir(os.path.dirname(fpath))

    torch.save(state, fpath)


class Logger(object):
    """
    Write console output to external text file.
    Code imported from https://github.com/Cysu/open-reid/blob/master/reid/utils/logging.py.
    """

    def __init__(self, fpath=None):
        self.console = sys.stdout
        self.file = None
        if fpath is not None:
            if not os.path.exists(os.path.dirname(fpath)):
                os.mkdir(os.path.dirname(fpath))

            self.file = open(fpath, 'w')

    def __del__(self):
        self.close()

    def __enter__(self):
        pass

    def __exit__(self, *args):
        self.close()

    def write(self, msg):
        self.console.write(msg)
        if self.file is not None:
            self.file.write(msg)

    def flush(self):
        self.console.flush()
        if self.file is not None:
            self.file.flush()
            os.fsync(self.file.fileno())

    def close(self):
        self.console.close()
        if self.file is not None:
            self.file.close()
