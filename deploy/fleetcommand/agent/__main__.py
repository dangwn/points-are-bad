from .launcher import Launcher
import argparse
import sys

args = sys.argv[1:]

parser = argparse.ArgumentParser()
parser.add_argument('-f', '--file-name')

launcher = Launcher(
    command = args[0],
    args = parser.parse_args(args[1:])
)

launcher.launch()