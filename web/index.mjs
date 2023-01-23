import express from 'express';
import { handler as ssrHandler } from './dist/server/entry.mjs';
import stream from 'node:stream';
import httpMocks from 'node-mocks-http';

var ws = new stream.Writable({
  write: function(chunk, encoding, next) {
    console.log(chunk.toString());
    next();
  }
});

const app = express();
app.use(express.static('dist/client/'))
app.use(async (req, res, next) => {
  await ssrHandler(req, res, next);
});

app.listen(8080);
