import { describe, expect, test } from "bun:test";
import { ChunkedReader, HttpRequestParser } from "../src/request";

describe("HTTP Request Parser", () => {
  test("parses GET request line with chunked reading", async () => {
    const testData = "GET /coffee HTTP/1.1\r\nHost: localhost\r\n\r\n";
    const parser = new HttpRequestParser();
    const reader = new ChunkedReader(testData, 3, 10);

    const promise = new Promise((resolve) => {
      parser.on("request-line", resolve);
    });

    await reader.readToParser(parser);
    const requestLine = await promise;

    expect(requestLine).toEqual({
      method: "GET",
      requestTarget: "/coffee",
      httpVersion: "1.1",
    });
  });
});
