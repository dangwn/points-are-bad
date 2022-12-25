from .utils import create_values_yaml
import subprocess
from argparse import Namespace
import os

DEFAULT_PAB_FILENAME = 'pab.yml'

class Launcher:
    def __init__(self, command: str, args: Namespace):
        self.command = command

        if not args.file_name:
            args.file_name = DEFAULT_PAB_FILENAME

        self.pab_values = create_values_yaml(
            pab_yaml_path=args.file_name,
            values_base_yaml_path=os.path.join(os.getcwd(),'values_base.yml'),
            values_yaml_path=os.path.join(os.getcwd(),'helm_charts','points-are-bad','values.yaml'),
            write_yaml=True
        )

    def run(self):
        print(subprocess.run([
            'helm','install',
            self.pab_values['global']['appName'],
            os.path.join(os.getcwd(),'helm_charts','points-are-bad'),
            '--namespace', self.pab_values['global']['appName'],
            '--create-namespace'
        ], stdout=subprocess.PIPE).stdout)

    def launch(self):
        if self.command == 'run':
            self.run()
