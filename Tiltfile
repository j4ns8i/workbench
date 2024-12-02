### Build steps

load('ext://podman', 'podman_build_with_restart')

def build_api():
    podman_build_with_restart(
        'workbench/api',
        '.',
        'fastapi dev api/app.py',
        extra_flags=['-f', 'api/build/Dockerfile'],
        live_update=[
            sync('./api/src', '/app/src'),
        ]
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

apply_common()
apply_api()
