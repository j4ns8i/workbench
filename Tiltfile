### Build steps

load('ext://restart_process', 'docker_build_with_restart')

def build_api():
    docker_build_with_restart(
        'workbench/api',
        './api',
        'fastapi run src/app.py',
        dockerfile='api/build/Dockerfile',
        live_update=[
            sync('./api/src', '/app/src'),
        ],
    )

build_api()

### Deploy steps

deploy_dir = 'deploy/local/workbench/templates/'

def apply_common():
    common_templates = ['serviceaccount.yaml']
    for p in common_templates:
        k8s_yaml(os.path.join(deploy_dir, p))

def apply_api():
    api_templates = ['api/service.yaml', 'api/deployment.yaml']
    for p in api_templates:
        k8s_yaml(os.path.join(deploy_dir, p))
    k8s_resource(workload='local-workbench-api', port_forwards=[8000])

apply_common()
apply_api()
