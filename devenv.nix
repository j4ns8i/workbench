{ pkgs, lib, config, inputs, ... }:
  let
    pkgs-unstable = import inputs.nixpkgs-unstable { system = pkgs.stdenv.system; };
  in {
    # https://devenv.sh/packages/
    packages = with pkgs; [
      git
      docker # still requires a docker daemon to be configured externally
      tilt
      kubectl
      kubectx
      k3d
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

    processes.api.exec = "fastapi dev api/app.py";

    tasks = {
      "build:api" = {
        exec = "docker build ./api -f api/build/Dockerfile -t workbench/api";
      };
      "helm:render:local" = {
        exec = "helm template local deploy/ --release-name --output-dir=deploy/";
      };
      "k3d:up" = {
        exec = "k3d cluster create workbench --registry-create workbench-registry";
        status = "k3d cluster get workbench &>/dev/null";
      };
      "k3d:down" = {
        exec = "k3d cluster delete workbench";
      };
    };

    # See full reference at https://devenv.sh/reference/options/
  }
