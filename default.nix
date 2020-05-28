{ buildGoModule
, nix-gitignore
}:

buildGoModule {
  pname = "archetype-go-axon";
  version = "0.0.1";
  src = nix-gitignore.gitignoreSource [] ./.;
  goPackagePath = "github.com/jeroenvm/archetype-go-axon";
  goDeps = ./deps.nix;
  modSha256 = "02apy9v62694p0gypk5lsm46giycmq0d6sjhql1azdr7ph8h9ggb";
}
