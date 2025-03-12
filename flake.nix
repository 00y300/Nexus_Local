{
  description = "Node.js development environment with latest Node.js and local bin path";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable"; # Use the latest nixpkgs
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

        nodejs = pkgs.nodejs; # Fetches the latest Node.js version from nixpkgs
      in
      rec {
        flakedPkgs = pkgs;

        # Enables use of `nix develop`
        devShell = pkgs.mkShell {

          hardeningDisable = [ "fortify" ];

          buildInputs = with pkgs; [
            nodejs
            eslint_d
            go
            gosimports
            gofumpt
            gotools
            delve
            gopls
            go-outline
            atac
            # golangci-lint

            tree
          ];

          # Add local node_modules/.bin to PATH
          shellHook = ''
            export PATH="$PWD/node_modules/.bin/:$PATH"
            export ATAC_KEY_BINDINGS="./vim.toml"
            echo "🐢 Node.js $(node -v) and npm $(npm -v) are ready to use!"
            echo "💡 Remember: Locally installed tools in node_modules/.bin are now available globally in this shell."
          '';
        };
      }
    );
}
