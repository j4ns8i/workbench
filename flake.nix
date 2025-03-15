{
  inputs = {
    utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, utils }: utils.lib.eachDefaultSystem (system:
    let
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      devShell = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          python311
          uv

          git
          docker
          tilt
          go-task
          kubernetes-helm
          kubectl
          kubectx
          stern
          k3d
          redis
        ];
      };
    }
  );
}
