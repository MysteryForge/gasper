{
  inputs = {
    nixpkgs = { url = "github:NixOS/nixpkgs/nixos-unstable"; };
    systems.url = "github:nix-systems/default";
    devshell.inputs.nixpkgs.follows = "nixpkgs";
    devshell.url = "github:numtide/devshell";
    foundry.url = "github:shazow/foundry.nix/monthly"; # Use monthly branch for permanent releases
  };

  outputs = { self, nixpkgs, systems, ... }@inputs:
    let
      eachSystem = f:
        nixpkgs.lib.genAttrs (import systems) (system:
          let
            pkgs = import nixpkgs {
              inherit system;
              config = { allowUnfree = true; };
              overlays = [];
            };
            devshell = pkgs.callPackage inputs.devshell { inherit inputs; };
            foundry = pkgs.callPackage "${inputs.foundry}/foundry-bin" {};
          in
          f pkgs devshell foundry
        );
    in {

      devShells = eachSystem (pkgs: devshell: foundry: {
        default = devshell.mkShell {
          packages = [
            pkgs.delve
            pkgs.gcc
            pkgs.go_1_24
            pkgs.gotools
            pkgs.gopls
            pkgs.go-outline
            pkgs.gopkgs
            pkgs.godef
            pkgs.golangci-lint
            pkgs.go-tools
            pkgs.treefmt
            pkgs.influxdb
            pkgs.gosec
            pkgs.jq
            pkgs.solc
            pkgs.go-ethereum
            pkgs.gitAndTools.git-absorb
            pkgs.act
            pkgs.gnumake

            foundry
          ];
        };
      });
    };
}