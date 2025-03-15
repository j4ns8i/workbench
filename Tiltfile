load('ext://helm_resource', 'helm_resource', 'helm_repo')
load('ext://restart_process', 'docker_build_with_restart')

### Build steps

def build_api():
    docker_build_with_restart(
        ref='workbench/api',
        context='./api',
        entrypoint='fastapi run src/app.py',
        dockerfile='api/build/Dockerfile',
        live_update=[
            sync('./api/src', '/app/src'),
        ],
    )

build_api()

def build_product_store():
    name = 'product-store'
    local_resource(
        name='bin.' + name,
        cmd='CGO_ENABLED=0 GOOS=linux go build -o ./bin/' + name,
        deps=[
            './{}/{}'.format(name, basename) for basename in [
                'main.go',
                'models.go',
                'http.go',
                'go.mod',
                'go.sum',
            ]
        ],
        dir=name,
    )

    entrypoint = '/usr/local/bin/' + name
    docker_build_with_restart(
        ref='workbench/' + name,
        context=name,
        entrypoint=[entrypoint],
        dockerfile='{}/build/tilt.Dockerfile'.format(name),
        only=[
            './bin',
        ],
        live_update=[
            sync('./{n}/bin/{n}'.format(n=name), entrypoint),
        ],
    )

build_product_store()

### Deploy steps

def apply_helm_chart():
    chart = helm('deploy', name='workbench')
    k8s_yaml(chart)
    k8s_resource(workload='workbench-api', port_forwards=[8000])
    k8s_resource(workload='workbench-product-store', port_forwards=[8080])

apply_helm_chart()

### Redis

def apply_redis_chart():
    helm_repo('bitnami', 'https://charts.bitnami.com/bitnami')
    helm_resource('redis', 'bitnami/redis',
        namespace='default',
        flags=[
            '--values=./deploy/redis-values.yaml',
            '--version=20.10.1',
        ],
        resource_deps=['bitnami'],
        deps=['./deploy/redis-values.yaml'],
    )
    k8s_resource(workload='redis', port_forwards=[6379])

apply_redis_chart()
