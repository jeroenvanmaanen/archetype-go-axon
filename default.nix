{ buildGoModule
, nix-gitignore
}:

buildGoModule {
  pname = "archetype-nix-go";
  version = "0.0.1";
  src = nix-gitignore.gitignoreSource [] ./.;
  goPackagePath = "github.com/jeroenvm/archetype-nix-go";
  goDeps = ./deps.nix;
  modSha256 = "0sjjj9z1dhilhpc8pq4154czrb79z9cm044jvn75kxcjv6v5l2m5";
}
