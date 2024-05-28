// package: comm
// file: comm.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";

export class BasicRequest extends jspb.Message { 

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): BasicRequest.AsObject;
    static toObject(includeInstance: boolean, msg: BasicRequest): BasicRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: BasicRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): BasicRequest;
    static deserializeBinaryFromReader(message: BasicRequest, reader: jspb.BinaryReader): BasicRequest;
}

export namespace BasicRequest {
    export type AsObject = {
    }
}

export class BasicResponse extends jspb.Message { 
    getCode(): number;
    setCode(value: number): BasicResponse;
    getMessage(): string;
    setMessage(value: string): BasicResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): BasicResponse.AsObject;
    static toObject(includeInstance: boolean, msg: BasicResponse): BasicResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: BasicResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): BasicResponse;
    static deserializeBinaryFromReader(message: BasicResponse, reader: jspb.BinaryReader): BasicResponse;
}

export namespace BasicResponse {
    export type AsObject = {
        code: number,
        message: string,
    }
}

export class GetValueRequest extends jspb.Message { 
    getKey(): string;
    setKey(value: string): GetValueRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetValueRequest.AsObject;
    static toObject(includeInstance: boolean, msg: GetValueRequest): GetValueRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetValueRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetValueRequest;
    static deserializeBinaryFromReader(message: GetValueRequest, reader: jspb.BinaryReader): GetValueRequest;
}

export namespace GetValueRequest {
    export type AsObject = {
        key: string,
    }
}

export class GetValueResponse extends jspb.Message { 
    getCode(): number;
    setCode(value: number): GetValueResponse;
    getMessage(): string;
    setMessage(value: string): GetValueResponse;
    getValue(): string;
    setValue(value: string): GetValueResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetValueResponse.AsObject;
    static toObject(includeInstance: boolean, msg: GetValueResponse): GetValueResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetValueResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetValueResponse;
    static deserializeBinaryFromReader(message: GetValueResponse, reader: jspb.BinaryReader): GetValueResponse;
}

export namespace GetValueResponse {
    export type AsObject = {
        code: number,
        message: string,
        value: string,
    }
}

export class SetValueRequest extends jspb.Message { 
    getKey(): string;
    setKey(value: string): SetValueRequest;
    getValue(): string;
    setValue(value: string): SetValueRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SetValueRequest.AsObject;
    static toObject(includeInstance: boolean, msg: SetValueRequest): SetValueRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SetValueRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SetValueRequest;
    static deserializeBinaryFromReader(message: SetValueRequest, reader: jspb.BinaryReader): SetValueRequest;
}

export namespace SetValueRequest {
    export type AsObject = {
        key: string,
        value: string,
    }
}

export class SetValueResponse extends jspb.Message { 
    getCode(): number;
    setCode(value: number): SetValueResponse;
    getMessage(): string;
    setMessage(value: string): SetValueResponse;
    getValue(): string;
    setValue(value: string): SetValueResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SetValueResponse.AsObject;
    static toObject(includeInstance: boolean, msg: SetValueResponse): SetValueResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SetValueResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SetValueResponse;
    static deserializeBinaryFromReader(message: SetValueResponse, reader: jspb.BinaryReader): SetValueResponse;
}

export namespace SetValueResponse {
    export type AsObject = {
        code: number,
        message: string,
        value: string,
    }
}

export class StrlnValueRequest extends jspb.Message { 
    getKey(): string;
    setKey(value: string): StrlnValueRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): StrlnValueRequest.AsObject;
    static toObject(includeInstance: boolean, msg: StrlnValueRequest): StrlnValueRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: StrlnValueRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): StrlnValueRequest;
    static deserializeBinaryFromReader(message: StrlnValueRequest, reader: jspb.BinaryReader): StrlnValueRequest;
}

export namespace StrlnValueRequest {
    export type AsObject = {
        key: string,
    }
}

export class StrlnValueResponse extends jspb.Message { 
    getCode(): number;
    setCode(value: number): StrlnValueResponse;
    getMessage(): string;
    setMessage(value: string): StrlnValueResponse;
    getValue(): number;
    setValue(value: number): StrlnValueResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): StrlnValueResponse.AsObject;
    static toObject(includeInstance: boolean, msg: StrlnValueResponse): StrlnValueResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: StrlnValueResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): StrlnValueResponse;
    static deserializeBinaryFromReader(message: StrlnValueResponse, reader: jspb.BinaryReader): StrlnValueResponse;
}

export namespace StrlnValueResponse {
    export type AsObject = {
        code: number,
        message: string,
        value: number,
    }
}

export class DeleteValueRequest extends jspb.Message { 
    getKey(): string;
    setKey(value: string): DeleteValueRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DeleteValueRequest.AsObject;
    static toObject(includeInstance: boolean, msg: DeleteValueRequest): DeleteValueRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: DeleteValueRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DeleteValueRequest;
    static deserializeBinaryFromReader(message: DeleteValueRequest, reader: jspb.BinaryReader): DeleteValueRequest;
}

export namespace DeleteValueRequest {
    export type AsObject = {
        key: string,
    }
}

export class DeleteValueResponse extends jspb.Message { 
    getCode(): number;
    setCode(value: number): DeleteValueResponse;
    getMessage(): string;
    setMessage(value: string): DeleteValueResponse;
    getValue(): string;
    setValue(value: string): DeleteValueResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DeleteValueResponse.AsObject;
    static toObject(includeInstance: boolean, msg: DeleteValueResponse): DeleteValueResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: DeleteValueResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DeleteValueResponse;
    static deserializeBinaryFromReader(message: DeleteValueResponse, reader: jspb.BinaryReader): DeleteValueResponse;
}

export namespace DeleteValueResponse {
    export type AsObject = {
        code: number,
        message: string,
        value: string,
    }
}

export class AppendValueRequest extends jspb.Message { 
    getKey(): string;
    setKey(value: string): AppendValueRequest;
    getValue(): string;
    setValue(value: string): AppendValueRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AppendValueRequest.AsObject;
    static toObject(includeInstance: boolean, msg: AppendValueRequest): AppendValueRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: AppendValueRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): AppendValueRequest;
    static deserializeBinaryFromReader(message: AppendValueRequest, reader: jspb.BinaryReader): AppendValueRequest;
}

export namespace AppendValueRequest {
    export type AsObject = {
        key: string,
        value: string,
    }
}

export class AppendValueResponse extends jspb.Message { 
    getCode(): number;
    setCode(value: number): AppendValueResponse;
    getMessage(): string;
    setMessage(value: string): AppendValueResponse;
    getValue(): string;
    setValue(value: string): AppendValueResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AppendValueResponse.AsObject;
    static toObject(includeInstance: boolean, msg: AppendValueResponse): AppendValueResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: AppendValueResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): AppendValueResponse;
    static deserializeBinaryFromReader(message: AppendValueResponse, reader: jspb.BinaryReader): AppendValueResponse;
}

export namespace AppendValueResponse {
    export type AsObject = {
        code: number,
        message: string,
        value: string,
    }
}

export class Entry extends jspb.Message { 
    getTerm(): number;
    setTerm(value: number): Entry;
    getKey(): string;
    setKey(value: string): Entry;
    getValue(): string;
    setValue(value: string): Entry;
    getIsconfigentry(): boolean;
    setIsconfigentry(value: boolean): Entry;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Entry.AsObject;
    static toObject(includeInstance: boolean, msg: Entry): Entry.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Entry, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Entry;
    static deserializeBinaryFromReader(message: Entry, reader: jspb.BinaryReader): Entry;
}

export namespace Entry {
    export type AsObject = {
        term: number,
        key: string,
        value: string,
        isconfigentry: boolean,
    }
}

export class AppendEntriesRequest extends jspb.Message { 
    getTerm(): number;
    setTerm(value: number): AppendEntriesRequest;
    getLeaderid(): number;
    setLeaderid(value: number): AppendEntriesRequest;
    getPrevlogindex(): number;
    setPrevlogindex(value: number): AppendEntriesRequest;
    getPrevlogterm(): number;
    setPrevlogterm(value: number): AppendEntriesRequest;
    clearEntriesList(): void;
    getEntriesList(): Array<Entry>;
    setEntriesList(value: Array<Entry>): AppendEntriesRequest;
    addEntries(value?: Entry, index?: number): Entry;
    getLeadercommit(): number;
    setLeadercommit(value: number): AppendEntriesRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AppendEntriesRequest.AsObject;
    static toObject(includeInstance: boolean, msg: AppendEntriesRequest): AppendEntriesRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: AppendEntriesRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): AppendEntriesRequest;
    static deserializeBinaryFromReader(message: AppendEntriesRequest, reader: jspb.BinaryReader): AppendEntriesRequest;
}

export namespace AppendEntriesRequest {
    export type AsObject = {
        term: number,
        leaderid: number,
        prevlogindex: number,
        prevlogterm: number,
        entriesList: Array<Entry.AsObject>,
        leadercommit: number,
    }
}

export class AppendEntriesResponse extends jspb.Message { 
    getTerm(): number;
    setTerm(value: number): AppendEntriesResponse;
    getSuccess(): boolean;
    setSuccess(value: boolean): AppendEntriesResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AppendEntriesResponse.AsObject;
    static toObject(includeInstance: boolean, msg: AppendEntriesResponse): AppendEntriesResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: AppendEntriesResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): AppendEntriesResponse;
    static deserializeBinaryFromReader(message: AppendEntriesResponse, reader: jspb.BinaryReader): AppendEntriesResponse;
}

export namespace AppendEntriesResponse {
    export type AsObject = {
        term: number,
        success: boolean,
    }
}

export class RequestVoteRequest extends jspb.Message { 
    getTerm(): number;
    setTerm(value: number): RequestVoteRequest;
    getCandidateid(): number;
    setCandidateid(value: number): RequestVoteRequest;
    getLastlogindex(): number;
    setLastlogindex(value: number): RequestVoteRequest;
    getLastlogterm(): number;
    setLastlogterm(value: number): RequestVoteRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RequestVoteRequest.AsObject;
    static toObject(includeInstance: boolean, msg: RequestVoteRequest): RequestVoteRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RequestVoteRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RequestVoteRequest;
    static deserializeBinaryFromReader(message: RequestVoteRequest, reader: jspb.BinaryReader): RequestVoteRequest;
}

export namespace RequestVoteRequest {
    export type AsObject = {
        term: number,
        candidateid: number,
        lastlogindex: number,
        lastlogterm: number,
    }
}

export class RequestVoteResponse extends jspb.Message { 
    getTerm(): number;
    setTerm(value: number): RequestVoteResponse;
    getVotegranted(): boolean;
    setVotegranted(value: boolean): RequestVoteResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RequestVoteResponse.AsObject;
    static toObject(includeInstance: boolean, msg: RequestVoteResponse): RequestVoteResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RequestVoteResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RequestVoteResponse;
    static deserializeBinaryFromReader(message: RequestVoteResponse, reader: jspb.BinaryReader): RequestVoteResponse;
}

export namespace RequestVoteResponse {
    export type AsObject = {
        term: number,
        votegranted: boolean,
    }
}

export class ChangeMembershipRequest extends jspb.Message { 
    getClusteraddresses(): string;
    setClusteraddresses(value: string): ChangeMembershipRequest;
    getNewclusteraddresses(): string;
    setNewclusteraddresses(value: string): ChangeMembershipRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ChangeMembershipRequest.AsObject;
    static toObject(includeInstance: boolean, msg: ChangeMembershipRequest): ChangeMembershipRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ChangeMembershipRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ChangeMembershipRequest;
    static deserializeBinaryFromReader(message: ChangeMembershipRequest, reader: jspb.BinaryReader): ChangeMembershipRequest;
}

export namespace ChangeMembershipRequest {
    export type AsObject = {
        clusteraddresses: string,
        newclusteraddresses: string,
    }
}

export class ChangeMembershipResponse extends jspb.Message { 
    getCode(): number;
    setCode(value: number): ChangeMembershipResponse;
    getMessage(): string;
    setMessage(value: string): ChangeMembershipResponse;
    getValue(): string;
    setValue(value: string): ChangeMembershipResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ChangeMembershipResponse.AsObject;
    static toObject(includeInstance: boolean, msg: ChangeMembershipResponse): ChangeMembershipResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ChangeMembershipResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ChangeMembershipResponse;
    static deserializeBinaryFromReader(message: ChangeMembershipResponse, reader: jspb.BinaryReader): ChangeMembershipResponse;
}

export namespace ChangeMembershipResponse {
    export type AsObject = {
        code: number,
        message: string,
        value: string,
    }
}
