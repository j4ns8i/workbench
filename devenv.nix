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

    # See full reference at https://devenv.sh/reference/options/
  }
