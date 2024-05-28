import "./style.css";
import { CommServiceClient } from "./grpc/comm.client";

import {
  BasicRequest,
} from "./grpc/comm";
import {
  pingButton,
  getButton,
  setButton,
  strlenButton,
  delButton,
  appendButton,
  sendButton
} from "./binding.ts";
import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport";


// grpc
const transport = new GrpcWebFetchTransport({
  baseUrl: "https://10.1.78.161:60000/",
  format: "binary",
});

const client = new CommServiceClient(
  transport
);

// event binding
pingButton.onclick = async () => {
  const request = BasicRequest.create({ value: "ping" });

  try {
    const response = await client.ping(request);
    console.log(response);
  } catch (e) {
    console.log(e);
  }
}

getButton.onclick = async () => {
  // TODO
}

setButton.onclick = async () => {
  // TODO
}

strlenButton.onclick = async () => {
  // TODO
}

delButton.onclick = async () => {
  // TODO
}

appendButton.onclick = async () => {
  // TODO
}

sendButton.onclick = async () => {
  // TODO
}