{
  pkgs ? import <nixpkgs> {}
}:

let
  myPackages = [
    pkgs.gnumake
    # Add more packages here
    pkgs.go
    pkgs.python312
  ];
in

pkgs.buildEnv {
  name = "worker";
  paths = myPackages;
}
