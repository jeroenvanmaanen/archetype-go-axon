import (
  let rev = "b0bbacb52134a7e731e549f4c0a7a2a39ca6b481"; in
  fetchTarball {
    url = "https://github.com/NixOS/nixpkgs-channels/archive/${rev}.tar.gz";
    sha256 = "15ix4spjpdm6wni28camzjsmhz0gzk3cxhpsk035952plwdxhb67";
  }
) {}
