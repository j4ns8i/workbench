load('ext://helm_resource', 'helm_resource', 'helm_repo')
load('ext://restart_process', 'docker_build_with_restart')

### Build steps

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
    k8s_resource(workload='workbench-api', port_forwards=[8000])

apply_helm_chart()

### Redis

def apply_redis_chart():
    helm_repo('bitnami', 'https://charts.bitnami.com/bitnami')
    helm_resource('redis', 'bitnami/redis',
        namespace='default',
        flags=['--values=./deploy/redis-values.yaml'],
        resource_deps=['bitnami'],
        deps=['./deploy/redis-values.yaml'],
    )
    k8s_resource(workload='redis', port_forwards=[6379])

apply_redis_chart()
