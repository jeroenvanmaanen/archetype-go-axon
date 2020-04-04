{ buildGoModule
, nix-gitignore
}:

buildGoModule {
  pname = "archetype-go-axon";
  version = "0.0.1";
  src = nix-gitignore.gitignoreSource [] ./.;
  goPackagePath = "github.com/jeroenvm/archetype-go-axon";
  goDeps = ./deps.nix;
  modSha256 = "0fli1cnh6y5i8aww43i80kchd2sz1wmcd55hycr0mq0cnqpvzmzv";
}
