package main

var indexHtml = `<html>

<head>
    <meta charset="utf-8">
    <style type="text/css">
        body {
            margin: 0;
        }

        .console {
            padding: 4px 8px;
            background-color: black;
            overflow-x: auto;
        }

        .console .mark {
            display: inline;
            color: white;
            font-size: small;
        }

		.console .loading {
			opacity: 0;
		}

        .console #stdout {
            display: block;
            margin: 0;
            padding: 0 0 4 0;
            width: 100%;
            height: 500px;
            font-size: 14px;
            color: white;
            overflow-y: auto;
            font-family: monospace;
            border-width: 0 0 1 0;
            border-color: white;
            border-style: solid;
        }

        .console #stdout::-webkit-scrollbar {
            width: 0;
        }

        .console #stdin {
            padding: 4 0 100 20;
            margin: 0 0 0 -20;
            width: 100%;
            border-width: 0;
            background-color: transparent;
            font-size: 14px;
            color: white;
            font-family: monospace;
        }

        .console #stdin:focus {
            outline: none;
        }
    </style>
    <script src="wasm_exec.js"></script>
    <script>
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("main.wasm", { credentials: "same-origin" }), go.importObject).then((result) => {
			let mark = document.getElementsByClassName("mark");
			if (mark) mark[0].className = "mark";
            go.run(result.instance);
        }).catch((err) => {
			let stdout = document.getElementById("stdout");
			if (stdout) stdout.innerText = err.toString();
		});
    </script>
</head>

<body>
    <div class="console">
        <pre id="stdout"></pre>
        <span class="mark loading">&gt; </span><input id="stdin" autofocus />
    </div>
</body>

</html>
`
