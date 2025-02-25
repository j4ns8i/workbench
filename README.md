# workbench

This is a general-purpose repo for building toy applications and proving out
ideas.

## Projects

### `api`

This project is a simple Python fastapi server which uses Redis as a backend
and offers publishing and subscribing to "topics", optionally using filters.
This is implemented using [Redis Streams] and [server-sent events].

[Redis Streams]: https://redis.io/docs/latest/develop/data-types/streams/
[server-sent events]: https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events

Example of subscribing to a topic "asdf" and filtering where "kind" ==
"message":

```bash
$ curl -v "http://localhost:8000/api/v0/events?kind=message&topic=asdf"
```

Example of publishing a "message" to the topic "asdf":

```bash
$ curl -H 'content-type: application/json' -XPOST http://localhost:8000/api/v0/messages -d '{"topic": "asdf", "message": "hello world"}'
```

## Development

This repo uses [devenv] ([nix]) to manage dependencies. [Task] is used in place
of make as a general-purpose command runner.

[devenv]: https://devenv.sh/
[nix]: https://nixos.org/guides/how-nix-works/
[Task]: https://taskfile.dev/

A local development environment can be launched on a [k3d] Kubernetes cluster
and managed by [Tilt] using `task up-tilt`.

[k3d]: https://k3d.io/stable/
[Tilt]: https://docs.tilt.dev/

API tests can be run using [pytest] with `task test-api`.

[pytest]: https://docs.pytest.org/en/stable/
