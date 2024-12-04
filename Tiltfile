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

def apply_helm_chart():
    chart = helm('deploy', name='workbench')
    k8s_yaml(chart)

apply_helm_chart()
