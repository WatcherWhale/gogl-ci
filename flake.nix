{
  description = "A cli tool for getting insight into your gitlab pipelines";

  outputs =
    {
      self,
      nixpkgs,
      flake-parts,
      devenv,
      gomod2nix,
      ...
    }@inputs:
    flake-parts.lib.mkFlake { inherit inputs; } (
      { ... }:
      {
        imports = [
          flake-parts.flakeModules.easyOverlay
          devenv.flakeModule
        ];

        systems = nixpkgs.lib.systems.flakeExposed;

        perSystem =
          {
            config,
            pkgs,
            system,
            ...
          }:
          {
            overlayAttrs = {
              inherit (config.packages) gogl-ci;
            };

            packages.default = config.packages.gogl-ci;

            packages.gogl-ci = gomod2nix.legacyPackages."${system}".buildGoApplication {
              pname = "gogl-ci";
              version = "0.0.0";

              src = "${self}";
              pwd = "${self}";
              modules = "${self}/gomod2nix.toml";

              meta = {
                mainProgram = "gogl";
              };

              ldflags = [
                "-w"
                "-s"
              ];

              subPackages = [ "cmd/cli" ];

              postInstall = ''
                mv $out/bin/cli $out/bin/gogl
              '';

            };

            devenv.shells.default = {
              packages = with pkgs; [
                go
                gotools
                golangci-lint
                go-task
                yaegi
                gomod2nix.packages."${system}".default
              ];
            };
          };
      }
    );

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";

    flake-parts = {
      url = "github:hercules-ci/flake-parts";
    };

    devenv = {
      url = "github:cachix/devenv";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

}
