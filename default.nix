{ buildGoModule
, nix-gitignore
}:

buildGoModule {
  pname = "archetype-nix-go";
  version = "0.0.1";
  src = nix-gitignore.gitignoreSource [] ./.;
  goPackagePath = "github.com/jeroenvm/archetype-nix-go";
  goDeps = ./deps.nix;
  modSha256 = "0mrwici0hn8q73yfavbxjnkly35c46si0dwp13j3nb0brfl53mzy";
}
