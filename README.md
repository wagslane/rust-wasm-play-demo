# rust-wasm-play-demo

This is a repo with some *bare banes* to show the proof of concept for a WASM Rust playground. There is a more feature complete live version on [Boot.dev here](https://boot.dev/playground/rs).

## Dependencies

* rustup
* wasm32-unknown-unknown

```bash
curl -sS https://webi.sh/rustlang | sh
rustup target add wasm32-unknown-unknown
```

## Clone

```bash
git clone https://github.com/wagslane/rust-wasm-play-demo
cd rust-wasm-play-demo
```

## Build

```bash
go build
```

## Run

```bash
./rust-wasm-play-demo
```

## Open webpage

Open [http://localhost:5000/app/](http://localhost:5000/app/) in your browser.
