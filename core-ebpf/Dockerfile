# Use the official Rust image as a build environment
FROM rust:1.80-slim-bullseye AS builder

# Install required build dependencies for eBPF and cross-compilation
RUN apt-get update && \
    apt-get install -y build-essential pkg-config libssl-dev musl-tools clang llvm

# Add the Linux musl target to build a statically linked binary
RUN rustup target add x86_64-unknown-linux-musl

WORKDIR /usr/src/stunt-ebpf

# Copy the Cargo configuration
COPY Cargo.toml Cargo.lock ./

# Create a dummy main.rs to cache dependency compilation
RUN mkdir src && \
    echo "fn main() {}" > src/main.rs && \
    cargo build --release --target x86_64-unknown-linux-musl && \
    rm -rf src

# Now copy the actual StuntDouble Rust source code
COPY src ./src

# Build the final StuntDouble eBPF injector binary
RUN touch src/main.rs && \
    cargo build --release --target x86_64-unknown-linux-musl

# ---------------------------------------------------
# Stage 2: Create an ultra-lightweight scratch image
# ---------------------------------------------------
FROM scratch

# Copy the statically linked Rust binary from the builder stage
COPY --from=builder /usr/src/stunt-ebpf/target/x86_64-unknown-linux-musl/release/core-ebpf /core-ebpf

# Run the StuntDouble eBPF Engine
ENTRYPOINT ["/core-ebpf"]
