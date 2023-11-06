settings = read_json('tilt_option.json', default={})

# Configure Tilt registry from imported settings
default_registry(settings.get('default_registry'))

# Specify which Kubernetes object manages the in-development container image
k8s_yaml('deployment.yaml')

# Specify name of the in-development container image and its build directory
docker_build('salman1091/hello-world','.')

# Specify name of k8s resource for Tilt to be aware of, and which ports should
# be forwarded to local machine
k8s_resource('hello-world', port_forwards=9000)
allow_k8s_contexts('kubernetes-admin@kubernetes')
