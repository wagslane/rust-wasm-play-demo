# rust-wasm-play-demo

This is a bare-bones repo to show the proof of concept for a WASM Rust playground. There is a more feature complete live version on [Boot.dev here](https://boot.dev/playground/rs).

## 1. Install Dependencies

* rustup
* wasm32-unknown-unknown

```bash
curl -sS https://webi.sh/rustlang | sh
rustup target add wasm32-unknown-unknown
```

## 2. Clone

```bash
git clone https://github.com/wagslane/rust-wasm-play-demo
cd rust-wasm-play-demo
```

## 3. Build

```bash
go build
```

## 4. Run

```bash
./rust-wasm-play-demo
```

## 5. Wait

Wait a long time because Rust is *not* blazingly fast at compiling.

## 5. Open webpage

Open [http://localhost:5000/app/](http://localhost:5000/app/) in your browser.
