FROM rust:1.67.1

RUN apt-get update && apt upgrade -y && apt-get -y install clang gcc g++ zlib1g-dev libmpc-dev libmpfr-dev libgmp-dev git cmake g++-aarch64-linux-gnu libc6-dev-arm64-cross && apt-get clean

# Add macOS Rust target
RUN rustup target add x86_64-apple-darwin && rustup target add aarch64-apple-darwin

RUN git clone https://github.com/tpoechtrager/osxcross
WORKDIR osxcross
RUN wget -nc https://s3.dockerproject.org/darwin/v2/MacOSX10.10.sdk.tar.xz && mv MacOSX10.10.sdk.tar.xz tarballs/
RUN UNATTENDED=yes OSX_VERSION_MIN=10.7 ./build.sh && ./build_gcc.sh
ENV PATH="/osxcross/target/bin:$PATH"
ENV CC=o64-clang
ENV CXX=o64-clang++
ENV LIBZ_SYS_STATIC=1