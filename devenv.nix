{ pkgs, lib, config, inputs, ... }:
  let
    pkgs-unstable = import inputs.nixpkgs-unstable { system = pkgs.stdenv.system; };
  in {
    # https://devenv.sh/packages/
    packages = with pkgs; [
      git
      podman # still requires a podman service to be configured externally
      tilt
      kubectl
      kubectx
      kind
      kubernetes-helm
    ];

    languages.python = {
      enable = true;
      uv = {
        enable = true;
        sync.enable = true;
        package = pkgs-unstable.uv;
      };
      venv = {
        enable = true;
      };
    };

    processes.python-api.exec = "fastapi dev python-api/app.py";

    tasks = {
      "helm:render" = {
        exec = "helm template tilt deploy/tilt > deploy/tilt/out/tilt.yaml";
      };
      "kind:up" = {
        exec = "systemd-run --scope --user -p \"Delegate=yes\" kind create cluster --name workbench";
        status = "kind get clusters 2>/dev/null | grep -q '^workbench$'";
      };
      "kind:down" = {
        exec = "kind delete cluster --name workbench";
      };
    };

    # See full reference at https://devenv.sh/reference/options/
  }
