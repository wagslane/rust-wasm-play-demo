const worker = new window.Worker(`/app/rs_worker.js`);

function useWorker(code, callback) {
    const promise = new Promise((resolve, reject) => {
        worker.onmessage = (event) => {
            if (event.data.done) {
                resolve();
                return;
            }
            if (event.data.error) {
                reject(event.data.error);
                return;
            }
            callback(event.data.message);
        };
    });
    worker.postMessage({ type: "EXEC_CODE", code });
    return promise;
}

async function compileRust(code) {
    const resp = await fetch(`/v1/compile`, {
        method: "POST",
        mode: "cors",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            code,
        }),
    });
    const handled = await handleWasmResponse(resp);
    return handled;
}

async function handleWasmResponse(response) {
    if (!response.ok) {
        const json = await response.json();
        if (typeof json.error !== "undefined") {
            throw json.error;
        }
        throw "Unknown error occured";
    }
    return await response.arrayBuffer();
}

async function clickRun() {
    const code = document.getElementById("code").value;
    const output = document.getElementById("output");

    console.log("Compiling...")
    output.innerHTML = "Compiling...";
    const wasm = await compileRust(code);

    console.log("Running...")
    output.innerHTML = "Running...";

    let count = 0;
    useWorker(wasm, (stdout) => {
        if (count === 0) {
            output.innerHTML = "";
            count++
        }
        output.innerHTML += stdout + "\n";
    })
}
