#!/bin/bash
#cargo install -f cross

# linux amd64 - OK
#cross build --release --target x86_64-unknown-linux-gnu
# windows x64 - OK
# cross build --release --target x86_64-pc-windows-gnu
# mac amd64 - OK - il me manque peut être des trucs
#cross build --release --target x86_64-apple-darwin

# linux amd64 - OK
#cross build --release --target aarch64-unknown-linux-gnu
# mac arm64 - NOK
#cross build --release --target aarch64-apple-darwin
# windows x64 - NOK
#cross build --release --target aarch64-pc-windows-msvc

# darwin x64

# darwin arm64
rust_version="1.72.0"
arch="aarch64-apple-darwin"
docker run --rm -v `pwd`:/app -w /usr/src/app rust:$rust_version sh -c "apt-get update && apt-get install -y g++-aarch64-linux-gnu libc6-dev-arm64-cross && cd /app && rustup target add $arch && rustup toolchain install stable-$arch && cargo build --release --target=$arch"
#mkdir -p dist/$arch && cp target/$arch/release/jeedom-status dist/$arch

#cargo build --release --target=x86_64-apple-darwin
#mkdir -p dist/darwin_amd64 && cp target/x86_64-apple-darwin/release/cep dist/darwin_amd64
