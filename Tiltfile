load('ext://helm_resource', 'helm_resource', 'helm_repo')

### Redis
def apply_redis_chart():
    helm_repo('bitnami', 'https://charts.bitnami.com/bitnami')
    helm_resource('redis', 'bitnami/redis',
        namespace='default',
        flags=[
            '--set=architecture=standalone',
            '--set=tls.authClients=false',
            '--set-json=master.disableCommands=[]',
            '--version=20.10.1',
        ],
        resource_deps=['bitnami'],
        deps=['./deploy/redis-values.yaml'],
    )
    k8s_resource(workload='redis', port_forwards=[6379])

apply_redis_chart()

### Component loading
load_dynamic('./product-store/Tiltfile')
