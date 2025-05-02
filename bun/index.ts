import chalk from "chalk";

// Bun.listen({
//   hostname: "localhost",
//   port: 42069,
//   socket: {
//     data: (s, data) => {
//       console.log(new TextDecoder().decode(data));
//     },
//   },
// });

Bun.connect({
  hostname: "localhost",
  port: 43000,
  socket: {
    open(s) {
      const rawHttp = [
        "GET / HTTP/1.1",
        "Host: localhost:42069",
        "User-Agent: BunClient",
        "",
        "",
      ].join("\r\n");

      s.write(rawHttp);
    },
    data(s, data) {
      console.log("---------------- 🚀 FROM CLIENT 🚀 ----------------");
      console.log("Server Response:\n", new TextDecoder().decode(data));
    },
  },
});

Bun.serve({
  port: 43000,
  fetch(req) {
    console.log("---------------- 🌐 FROM WEB SERVER 🌐 ----------------");
    console.log(chalk.black.bgGreen.bold(`Method: ${req.method} ${req.url}\n`));

    return new Response("Hello from Bun.serve!");
  },
});
