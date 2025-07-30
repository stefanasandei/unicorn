{
  pkgs ? import <nixpkgs> {}
}:

let
  myPackages = [
    pkgs.gnumake
    # Add more packages here
    %s
  ];
in

pkgs.buildEnv {
  name = "worker";
  paths = myPackages;
}
