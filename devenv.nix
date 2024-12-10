{ pkgs, lib, config, inputs, ... }:
  let
    pkgs-unstable = import inputs.nixpkgs-unstable { system = pkgs.stdenv.system; };
  in {

    # https://devenv.sh/packages/
    packages = with pkgs; [
      git
      docker # still requires a docker daemon to be configured externally
      tilt
      go-task
      kubernetes-helm
    ] ++ (with pkgs-unstable; [
      kubectl
      kubectx
      stern
      k3d
      redis
    ]);

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
