{ buildGoModule
, nix-gitignore
}:

buildGoModule {
  pname = "archetype-nix-go";
  version = "0.0.1";
  src = nix-gitignore.gitignoreSource [] ./.;
  goPackagePath = "github.com/jeroenvm/archetype-nix-go";
  goDeps = ./deps.nix;
  modSha256 = "096kffzh8nlw0xpgxdzz01an768sj76x09qmg4i7wga02kgi2f95";
}
