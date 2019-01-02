{ pkgs ? import ./pkgs.nix }: with pkgs;

buildGoPackage {
  name = "seed-ca";
  src = lib.cleanSource ./.;

  goPackagePath = "gitlab.com/transumption/unstable/seed-ca";
  doCheck = true;
}
