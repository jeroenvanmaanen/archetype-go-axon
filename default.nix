{ buildGoModule
, nix-gitignore
}:

buildGoModule {
  pname = "archetype-go-axon";
  version = "0.0.1";
  src = nix-gitignore.gitignoreSource [] ./.;
  goPackagePath = "github.com/jeroenvm/archetype-go-axon";
  goDeps = ./deps.nix;
  modSha256 = "0m7s6qk61m2sjssgrd8zzxnjkypn6gi279pk2jcqm0g0vc2r71hi";
}
