{
  description = "A cli tool for getting insight into your gitlab pipelines";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
    };
    devenv = {
      url = "github:cachix/devenv";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = {
    nixpkgs,
    flake-parts,
    devenv,
    ...
  } @ inputs: let
    version = "0.0.2";
    rev = "e052d8e881f30ea89fcffad150bdce60e2fdc6bd";
    hash = "sha256-h20+8EtuYL6mwyaqgIvYAxWnUfUbISVmcsWp/LfUVtA=";
    vendorHash = "sha256-H1pwnlwvegDZ20iag8wltaE7BJg4mPV7MloE0GLuQpM=";
  in
    flake-parts.lib.mkFlake {inherit inputs;} ({...}: {
      imports = [
        flake-parts.flakeModules.easyOverlay
        devenv.flakeModule
      ];

      systems = nixpkgs.lib.systems.flakeExposed;

      perSystem = {
        config,
        pkgs,
        ...
      }: {
        overlayAttrs = {
          inherit (config.packages) gogl-ci;
        };

        formatter = pkgs.alejandra;

        packages.default = config.packages.gogl-ci;

        packages.gogl-ci = nixpkgs.legacyPackages.x86_64-linux.buildGoModule {
          pname = "gogl-ci";
          version = "${version}";

          vendorHash = "${vendorHash}";
          proxyVendor = true;
          tags = ["netgo"];

          src = nixpkgs.legacyPackages.x86_64-linux.fetchgit {
            name = "gogl-ci";
            url = "https://github.com/WatcherWhale/gogl-ci.git";
            rev = "${rev}";
            hash = "${hash}";
          };

          ldflags = [
            "-w"
            "-s"
          ];

          meta = {
            descirption = "GoGL-CI: A cli tool for getting insight into your gitlab pipelines";
            mainProgram = "gogl";
          };

          subPackages = ["cmd/cli"];

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
            alejandra
            yaegi
          ];
        };
      };
    });
}
