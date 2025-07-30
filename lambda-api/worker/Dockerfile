FROM nixos/nix 

COPY . ./app

WORKDIR ./app

RUN nix-shell -p pkgs.python312 python3Packages.pyyaml --run "python3 runtimes/build_env.py"

RUN nix-env -i -f ./runtimes/default.nix

WORKDIR ./worker

RUN go get .

RUN go build .

CMD ./worker
