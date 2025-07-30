import os
import yaml


def generate_nix_file(template_path, runtime_path, output_path):
    # Load the template
    with open(template_path, 'r') as template_file:
        template_content = template_file.read()

    # Find and extract the placeholder (%s) position
    placeholder_position = template_content.find('%s')

    # Iterate over runtime files in the runtime directory
    nix_pkgs_list = []
    for runtime_file in os.listdir(runtime_path):
        if runtime_file.endswith('.yaml'):
            runtime_file_path = os.path.join(runtime_path, runtime_file)

            # Load the runtime YAML file
            with open(runtime_file_path, 'r') as runtime_file_content:
                runtime_data = yaml.safe_load(runtime_file_content)

            # Extract nix_pkgs information
            nix_pkgs = runtime_data.get('nix_pkgs', [])
            nix_pkgs_list.extend(nix_pkgs)

    # Insert the nix_pkgs information into the template
    nix_pkgs_str = '\n    '.join(nix_pkgs_list)
    new_content = template_content[:placeholder_position] + nix_pkgs_str + template_content[placeholder_position + 2:]

    # Write the new content to the output file
    with open(output_path, 'w') as output_file:
        output_file.write(new_content)


if __name__ == "__main__":
    # Specify the paths
    template_path = './runtimes/template.nix'
    runtime_path = './runtimes/'
    output_path = './runtimes/default.nix'

    # Generate the default.nix file
    generate_nix_file(template_path, runtime_path, output_path)

    print(f"Generated {output_path}")
