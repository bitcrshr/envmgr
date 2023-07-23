SHELL := /bin/bash

.PHONY: update-bin
update-bin:
	cargo build --release && \
	cd .. && \
	cp ./target/release/envmgr-protoc ./rust-proto-compile && \
	chmod +x ./rust-proto-compile && \
	echo "Updated ./rust-proto-compile"