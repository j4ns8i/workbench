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
      go-task
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

  }
