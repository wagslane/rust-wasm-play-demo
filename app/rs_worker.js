// keep a WebAssembly memory reference for `readString`
let memory;

// read a null terminated c string at a wasm memory buffer index
function readString(ptr) {
  const view = new Uint8Array(memory.buffer);

  // find the end of the string (null)
  let end = ptr;
  while (view[end]) ++end;

  // `subarray` uses the same underlying ArrayBuffer as the view
  const buf = new Uint8Array(view.subarray(ptr, end));
  const str = new TextDecoder().decode(buf); // (utf-8 by default)

  return str;
}

addEventListener(
  "message",
  async (e) => {
    console.log("Received compiled wasm...");

    const importObj = {
      env: {
        console_log: (line) => {
          postMessage({
            message: readString(line),
          });
        },
        console_warn: (line) => {
          postMessage({
            message: readString(line),
          });
        },
        console_error: (line) => {
          postMessage({
            error: readString(line),
          });
        },
      },
    };

    const obj = await WebAssembly.instantiate(e.data.code, importObj);
    memory = obj.instance.exports.memory;
    obj.instance.exports.lib();

    postMessage({
      done: true,
    });
    console.log("Done running code!");
  },
  false
);
