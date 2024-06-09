import "./style.css";
import { CommServiceClient } from "./grpc/comm.client";

import {
  AppendValueRequest,
  BasicRequest, ChangeMembershipRequest, DeleteValueRequest, GetValueRequest, SetValueRequest,
  StrlnValueRequest,
} from "./grpc/comm";
import {
  pingButton,
  getButton,
  setButton,
  strlenButton,
  delButton,
  appendButton,
  getInput, setKeyInput, setValueInput, addressInput,
  strlenInput,
  appendKeyInput,
  appendValueInput,
  logSection,
  delInput,
  changeMembershipButton,
  changeMembershipInput
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

// Update log functionality
function updateLog(message: string) {
  const logs = logSection;
  if (logs) {
    const newLogEntry = document.createElement('div');
    newLogEntry.textContent = `${new Date().toLocaleTimeString()} - ${addressInput.value} - ${message}`;
    logs.insertBefore(newLogEntry, logSection.firstChild);
  }
}

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
    updateLog(response.response.message)
  } catch (e) {
    console.log(e);
    updateLog("An error occured when trying to PING, check console")
  }
}

getButton.onclick = async () => {
  const request = GetValueRequest.create({
    key: getInput.value
  });

  try {
    const response = await client.getValue(request);
    console.log(response)
    if (getInput.value && response.response.value) {
      updateLog(response.response.message + " " + getInput.value + " = " + response.response.value)
    } else if (!getInput.value) {
      updateLog(response.response.message + ", but no key was used")
    } else {
      updateLog(response.response.message + " " + getInput.value + ", but no value was found")
    }
  } catch (e) {
    console.log(e);
    updateLog("An error occured when trying to GET value, check console")
  }
}

setButton.onclick = async () => {
  const request = SetValueRequest.create({
    key: setKeyInput.value,
    value: setValueInput.value
  });

  try {
    const response = await client.setValue(request);
    if (setKeyInput.value && setValueInput.value) {
      updateLog(response.response.message + " " + setKeyInput.value + " = " + setValueInput.value)
    } else if (!setKeyInput.value) {
      updateLog(response.response.message + ", but no key was used")
    } else {
      updateLog(response.response.message + ", but no value was used")
    }
  } catch (e) {
    console.log(e);
    updateLog("An error occured when trying to SET value, check console")
  }
}

strlenButton.onclick = async () => {
  const request = StrlnValueRequest.create({
    key: strlenInput.value
  })

  try {
    const response = await client.strlnValue(request)
    if (strlenInput.value) {
      updateLog(response.response.message + " length of " + strlenInput.value + " = " + response.response.value)
    } else {
      updateLog(response.response.message + ", but no key was used")
    }
  } catch (e) {
    console.log(e)
    updateLog("An error occured when trying to STRLEN key, check console")
  }
}

delButton.onclick = async () => {
  const request = DeleteValueRequest.create({
    key: delInput.value
  })

  try {
    const response = await client.deleteValue(request)
    if (delInput.value) {
      updateLog(response.response.message + ", deleted "+ delInput.value)
    } else {
      updateLog(response.response.message + ", but no key was used")
    }
  } catch (e) {
    console.log(e)
    updateLog("An error occured when trying to DELETE key, check console")
  }
}

appendButton.onclick = async () => {
  const request = AppendValueRequest.create({
    key: appendKeyInput.value,
    value: appendValueInput.value,
  })

  try {
    const response = await client.appendValue(request)
    if (appendKeyInput.value && appendValueInput.value) {
      updateLog(response.response.message + " " + appendKeyInput.value)
    } else if (!appendKeyInput.value) {
      updateLog(response.response.message + ", but no key was used")
    } else {
      updateLog(response.response.message + ", but no value was used")
    }
  } catch (e) {
    console.log(e)
    updateLog("An error occured when trying to APPEND value, check console")
  }
}

changeMembershipButton.onclick = async () => {
  const request = ChangeMembershipRequest.create({
    newClusterAddresses: changeMembershipInput.value,
  })

  try {
    const response = await client.changeMembership(request)
    
    updateLog(response.response.message)
  } catch (e) {
    console.log(e)
    updateLog("An error occured when trying to CHANGE MEMBERSHIP, check console")
  }
}