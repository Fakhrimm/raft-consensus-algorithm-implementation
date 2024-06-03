import "./style.css";
import { CommServiceClient } from "./grpc/comm.client";

import {
  BasicRequest, GetValueRequest, SetValueRequest,
} from "./grpc/comm";
import {
  pingButton,
  getButton,
  setButton,
  strlenButton,
  delButton,
  appendButton,
  sendButton, getInput, setKeyInput, setValueInput, addressInput
} from "./binding.ts";
import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport";


// grpc
let transport = new GrpcWebFetchTransport({
  baseUrl: "http://" + addressInput.value || "0.0.0.0:60000",
  format: "binary",
});

let client = new CommServiceClient(
  transport
);

// event binding
addressInput.onchange = () => {
  transport = new GrpcWebFetchTransport({
    baseUrl: "http://" + addressInput.value,
    format: "binary",
  });

  client = new CommServiceClient(
    transport
  );
}

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
  const request = GetValueRequest.create({
    key: getInput.value
  });

  try {
    const response = await client.getValue(request);
    console.log(response);
  } catch (e) {
    console.log(e);
  }
}

setButton.onclick = async () => {
  const request = SetValueRequest.create({
    key: setKeyInput.value,
    value: setValueInput.value
  });

  try {
    const response = await client.setValue(request);
    console.log(response);
  } catch (e) {
    console.log(e);
  }
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