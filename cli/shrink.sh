#!/bin/sh
WIN=dist/kvs_windows_amd64/kvs.exe
NIX=dist/kvs_linux_amd64/kvs
MAC=dist/kvs_darwin_amd64/kvs

if [ -x "$(command -v upx)" ]; then
    upx -9 ${WIN}
    
    mv ${NIX} ${NIX%.*}.elf
    upx -9 -o ${NIX} ${NIX%.*}.elf
    rm ${NIX%.*}.elf
        
    mv ${MAC} ${MAC%.*}.macho
    upx -9 -o ${MAC} ${MAC%.*}.macho
    rm ${MAC%.*}.macho
fi