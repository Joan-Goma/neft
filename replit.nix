{ pkgs }: {
    deps = [
        pkgs.zip
         pkgs.gotools
         pkgs.go_1_17
        pkgs.gopls
    ];
}