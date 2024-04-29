# default.nix

{ pkgs ? import <nixpkgs> {} }:

pkgs.python3.withPackages (ps: with ps; [
  # Python packages
  flask
  gunicorn
  # Add other dependencies as needed
])
