{ buildGoModule
, nix-gitignore
}:

buildGoModule {
  pname = "archetype-nix-go";
  version = "0.0.1";
  src = nix-gitignore.gitignoreSource [] ./.;
  goPackagePath = "github.com/jeroenvm/archetype-nix-go";
  goDeps = ./deps.nix;
  modSha256 = "12c4pmqg1jlhfiv73yfs5q40i9fzgnf5447jmfcysbiyzrj04mk3";
}
