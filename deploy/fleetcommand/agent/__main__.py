from .utils import create_values_yaml

yml = create_values_yaml(
    pab_yaml_path = 'pab_master.yml',
    values_base_yaml_path = 'values_base.yml',
    values_yaml_path = 'helm_charts/points-are-bad/values.yaml',
    write_yaml = True
)
print(yml)