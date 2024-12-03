{ pkgs, lib, config, inputs, ... }:
  let
    pkgs-unstable = import inputs.nixpkgs-unstable { system = pkgs.stdenv.system; };
    targets = {
      "build:api-image" = {
        exec = "docker build ./api -f api/build/Dockerfile -t workbench/api";
      };
      "build:helm" = {
        exec = "helm template local deploy/ --release-name --output-dir=deploy/";
        before = [ "up:k3d" ];
      };
      "up:k3d" = {
        exec = "k3d cluster create workbench --registry-create workbench-registry";
        status = "k3d cluster get workbench &>/dev/null";
        before = [ "up:tilt" ];
      };
      "down:k3d" = {
        exec = "k3d cluster delete workbench";
      };
      "up:tilt" = {
        exec = "tilt up";
        status = "tilt status";
      };
    };
    scripts = lib.mapAttrs' (name: _: {
      name = "dtr-${name}";
      value = {
        exec = "devenv tasks run ${name}";
      };
    }) targets;
    tasks = lib.mapAttrs (name: target: {
      exec = target.exec;
    } // lib.optionalAttrs (target ? status) {
      inherit (target) status;
    }) targets;
  in {
    inherit tasks scripts;

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

    # See full reference at https://devenv.sh/reference/options/
  }
