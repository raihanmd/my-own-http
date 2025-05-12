import * as net from "node:net";
import { EventEmitter } from "node:stream";

export class HttpRequestParser extends EventEmitter {
  private buffer: Buffer;
  private state: string;
  private initialBufferSize: number;

  constructor() {
    super();
    this.buffer = Buffer.alloc(0);
    this.state = "initialized";
    this.initialBufferSize = 8;
  }

  feed(data: Buffer<ArrayBufferLike>) {
    this.buffer = Buffer.concat([this.buffer, data]);

    this._parse();
  }

  _parse() {
    switch (this.state) {
      case "initialized":
        this._parseRequestLine();
        break;
      case "done":
        this.emit("error", new Error("Trying to read data in a done state"));
        break;
      default:
        this.emit("error", new Error("Unknown state"));
    }
  }

  _parseRequestLine() {
    // Find the end of the request line
    const endIndex = this.buffer.indexOf("\r\n");
    if (endIndex === -1) {
      return; // Need more data
    }

    // Extract and parse the request line
    const line = this.buffer.subarray(0, endIndex).toString();
    const parts = line.split(" ");

    if (parts.length !== 3) {
      this.emit("error", new Error("Invalid request line format"));
      return;
    }

    // Emit the parsed request
    this.emit("request-line", {
      method: parts[0],
      requestTarget: parts[1],
      httpVersion: parts[2]!.replace("HTTP/", ""),
    });

    // Remove consumed data from buffer
    this.buffer = this.buffer.subarray(endIndex + 2);
    this.state = "done";
  }
}

export class ChunkedReader {
  private pos: number;
  private data: Buffer;
  private bytesPerRead: number;
  private delayMs: number;

  constructor(data: string, bytesPerRead = 1, delayMs = 100) {
    this.data = Buffer.from(data);
    this.bytesPerRead = bytesPerRead;
    this.delayMs = delayMs;
    this.pos = 0;
  }

  async readToParser(parser: HttpRequestParser) {
    while (this.pos < this.data.length) {
      const end = Math.min(this.pos + this.bytesPerRead, this.data.length);
      const chunk = this.data.subarray(this.pos, end);
      this.pos += chunk.length;

      parser.feed(chunk);

      // Simulate network lag
      await new Promise((resolve) => setTimeout(resolve, this.delayMs));
    }
  }
}

const server = net.createServer((socket) => {
  socket.on("data", (data) => {
    console.log(new TextDecoder().decode(data));
  });
});

server.listen(42069, () => {
  console.log("🚀 Socket ready!");
});
