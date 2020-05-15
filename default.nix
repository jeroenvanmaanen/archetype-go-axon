{ buildGoModule
, nix-gitignore
}:

buildGoModule {
  pname = "archetype-go-axon";
  version = "0.0.1";
  src = nix-gitignore.gitignoreSource [] ./.;
  goPackagePath = "github.com/jeroenvm/archetype-go-axon";
  goDeps = ./deps.nix;
  modSha256 = "1c6j9fs95nghnmlawc1gyvzrkgb6kra5vsnnx14yaj627w2ip3mj";
}
