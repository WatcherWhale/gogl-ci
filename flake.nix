{
  description = "A cli tool for getting insight into your gitlab pipelines";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      ...
    }:
    let
      version = "0.0.1";
      rev = "0c3b049ee6837bdebfda18dc670999a20f9ac7cd";
      hash = "sha256-LxWrjZfGjcN10qg5XMWaCAF42KTUS5ziNRc+6cIeHVs=";
      vendorHash = "sha256-b+BT/IZ6rvXsLokhBWruoA71N3hUsbuipD4V9D1ypEQ=";
    in
    {
      devShells = flake-utils.lib.eachDefaultSystemPassThrough (
        system:
        let
          pkgs = import nixpkgs {
            system = system;
          };
        in
        {
          "${system}".default = pkgs.mkShell {
            packages = with pkgs; [
              go
              golangci-lint
              go-task
              yaegi
            ];
          };
        }
      );

      packages = flake-utils.lib.eachDefaultSystemPassThrough (system: {
        "${system}" = {
          default = self.packages."${system}".gogl-ci;
          gogl-ci = nixpkgs.legacyPackages.x86_64-linux.buildGoModule {
            pname = "gogl-ci";
            version = "${version}";

            vendorHash = "${vendorHash}";
            proxyVendor = true;
            tags = [ "netgo" ];

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

            subPackages = [ "cmd/cli" ];

            postInstall = ''
              mv $out/bin/cli $out/bin/gogl
            '';
          };
        };
      });
    };
}
