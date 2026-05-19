{
  description = "Nix flake for rtutils_lib";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };

        goModulesHash = "sha256-ckMi/IIx+OOU6VX7WoqU6YjIYKdIAWGBtdnqeWWhox8=";

        rtutilsTestSuite = pkgs.buildGoModule {
          pname = "rtutils-test-suite";
          version = "0.1.0";
          src = ./.;
          vendorHash = goModulesHash;
          subPackages = [ "cmd/rtutils-test-suite" ];
          ldflags = [
            "-s"
            "-w"
          ];

          meta = with pkgs.lib; {
            description = "Request Tracker API utilities test suite CLI";
            license = licenses.mit;
            platforms = platforms.unix;
          };
        };

        rtutilsLibCheck = pkgs.buildGoModule {
          pname = "rtutils-lib-check";
          version = "0.1.0";
          src = ./.;
          vendorHash = goModulesHash;

          # Build all library packages to verify the module compiles.
          buildPhase = ''
            runHook preBuild
            go test -run=^$ ./...
            runHook postBuild
          '';

          checkPhase = ''
            runHook preCheck
            go test ./...
            runHook postCheck
          '';

          installPhase = ''
            runHook preInstall
            mkdir -p "$out"
            cp -r . "$out/src"
            runHook postInstall
          '';
        };
      in
      {
        packages.default = rtutilsTestSuite;
        packages.rtutils-test-suite = rtutilsTestSuite;
        packages.rtutils-lib-check = rtutilsLibCheck;

        checks.default = rtutilsLibCheck;

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            gotools
          ];
        };
      }
    );
}
