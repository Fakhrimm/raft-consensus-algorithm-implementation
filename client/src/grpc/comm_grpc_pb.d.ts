// package: comm
// file: comm.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "grpc";
import * as comm_pb from "./comm_pb";

interface ICommServiceService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    ping: ICommServiceService_IPing;
    stop: ICommServiceService_IStop;
    getValue: ICommServiceService_IGetValue;
    setValue: ICommServiceService_ISetValue;
    strlnValue: ICommServiceService_IStrlnValue;
    deleteValue: ICommServiceService_IDeleteValue;
    appendValue: ICommServiceService_IAppendValue;
    changeMembership: ICommServiceService_IChangeMembership;
    appendEntries: ICommServiceService_IAppendEntries;
    requestVote: ICommServiceService_IRequestVote;
}

interface ICommServiceService_IPing extends grpc.MethodDefinition<comm_pb.BasicRequest, comm_pb.BasicResponse> {
    path: "/comm.CommService/Ping";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<comm_pb.BasicRequest>;
    requestDeserialize: grpc.deserialize<comm_pb.BasicRequest>;
    responseSerialize: grpc.serialize<comm_pb.BasicResponse>;
    responseDeserialize: grpc.deserialize<comm_pb.BasicResponse>;
}
interface ICommServiceService_IStop extends grpc.MethodDefinition<comm_pb.BasicRequest, comm_pb.BasicResponse> {
    path: "/comm.CommService/Stop";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<comm_pb.BasicRequest>;
    requestDeserialize: grpc.deserialize<comm_pb.BasicRequest>;
    responseSerialize: grpc.serialize<comm_pb.BasicResponse>;
    responseDeserialize: grpc.deserialize<comm_pb.BasicResponse>;
}
interface ICommServiceService_IGetValue extends grpc.MethodDefinition<comm_pb.GetValueRequest, comm_pb.GetValueResponse> {
    path: "/comm.CommService/GetValue";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<comm_pb.GetValueRequest>;
    requestDeserialize: grpc.deserialize<comm_pb.GetValueRequest>;
    responseSerialize: grpc.serialize<comm_pb.GetValueResponse>;
    responseDeserialize: grpc.deserialize<comm_pb.GetValueResponse>;
}
interface ICommServiceService_ISetValue extends grpc.MethodDefinition<comm_pb.SetValueRequest, comm_pb.SetValueResponse> {
    path: "/comm.CommService/SetValue";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<comm_pb.SetValueRequest>;
    requestDeserialize: grpc.deserialize<comm_pb.SetValueRequest>;
    responseSerialize: grpc.serialize<comm_pb.SetValueResponse>;
    responseDeserialize: grpc.deserialize<comm_pb.SetValueResponse>;
}
interface ICommServiceService_IStrlnValue extends grpc.MethodDefinition<comm_pb.StrlnValueRequest, comm_pb.StrlnValueResponse> {
    path: "/comm.CommService/StrlnValue";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<comm_pb.StrlnValueRequest>;
    requestDeserialize: grpc.deserialize<comm_pb.StrlnValueRequest>;
    responseSerialize: grpc.serialize<comm_pb.StrlnValueResponse>;
    responseDeserialize: grpc.deserialize<comm_pb.StrlnValueResponse>;
}
interface ICommServiceService_IDeleteValue extends grpc.MethodDefinition<comm_pb.DeleteValueRequest, comm_pb.DeleteValueResponse> {
    path: "/comm.CommService/DeleteValue";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<comm_pb.DeleteValueRequest>;
    requestDeserialize: grpc.deserialize<comm_pb.DeleteValueRequest>;
    responseSerialize: grpc.serialize<comm_pb.DeleteValueResponse>;
    responseDeserialize: grpc.deserialize<comm_pb.DeleteValueResponse>;
}
interface ICommServiceService_IAppendValue extends grpc.MethodDefinition<comm_pb.AppendValueRequest, comm_pb.AppendValueResponse> {
    path: "/comm.CommService/AppendValue";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<comm_pb.AppendValueRequest>;
    requestDeserialize: grpc.deserialize<comm_pb.AppendValueRequest>;
    responseSerialize: grpc.serialize<comm_pb.AppendValueResponse>;
    responseDeserialize: grpc.deserialize<comm_pb.AppendValueResponse>;
}
interface ICommServiceService_IChangeMembership extends grpc.MethodDefinition<comm_pb.ChangeMembershipRequest, comm_pb.ChangeMembershipResponse> {
    path: "/comm.CommService/ChangeMembership";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<comm_pb.ChangeMembershipRequest>;
    requestDeserialize: grpc.deserialize<comm_pb.ChangeMembershipRequest>;
    responseSerialize: grpc.serialize<comm_pb.ChangeMembershipResponse>;
    responseDeserialize: grpc.deserialize<comm_pb.ChangeMembershipResponse>;
}
interface ICommServiceService_IAppendEntries extends grpc.MethodDefinition<comm_pb.AppendEntriesRequest, comm_pb.AppendEntriesResponse> {
    path: "/comm.CommService/AppendEntries";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<comm_pb.AppendEntriesRequest>;
    requestDeserialize: grpc.deserialize<comm_pb.AppendEntriesRequest>;
    responseSerialize: grpc.serialize<comm_pb.AppendEntriesResponse>;
    responseDeserialize: grpc.deserialize<comm_pb.AppendEntriesResponse>;
}
interface ICommServiceService_IRequestVote extends grpc.MethodDefinition<comm_pb.RequestVoteRequest, comm_pb.RequestVoteResponse> {
    path: "/comm.CommService/RequestVote";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<comm_pb.RequestVoteRequest>;
    requestDeserialize: grpc.deserialize<comm_pb.RequestVoteRequest>;
    responseSerialize: grpc.serialize<comm_pb.RequestVoteResponse>;
    responseDeserialize: grpc.deserialize<comm_pb.RequestVoteResponse>;
}

export const CommServiceService: ICommServiceService;

export interface ICommServiceServer {
    ping: grpc.handleUnaryCall<comm_pb.BasicRequest, comm_pb.BasicResponse>;
    stop: grpc.handleUnaryCall<comm_pb.BasicRequest, comm_pb.BasicResponse>;
    getValue: grpc.handleUnaryCall<comm_pb.GetValueRequest, comm_pb.GetValueResponse>;
    setValue: grpc.handleUnaryCall<comm_pb.SetValueRequest, comm_pb.SetValueResponse>;
    strlnValue: grpc.handleUnaryCall<comm_pb.StrlnValueRequest, comm_pb.StrlnValueResponse>;
    deleteValue: grpc.handleUnaryCall<comm_pb.DeleteValueRequest, comm_pb.DeleteValueResponse>;
    appendValue: grpc.handleUnaryCall<comm_pb.AppendValueRequest, comm_pb.AppendValueResponse>;
    changeMembership: grpc.handleUnaryCall<comm_pb.ChangeMembershipRequest, comm_pb.ChangeMembershipResponse>;
    appendEntries: grpc.handleUnaryCall<comm_pb.AppendEntriesRequest, comm_pb.AppendEntriesResponse>;
    requestVote: grpc.handleUnaryCall<comm_pb.RequestVoteRequest, comm_pb.RequestVoteResponse>;
}

export interface ICommServiceClient {
    ping(request: comm_pb.BasicRequest, callback: (error: grpc.ServiceError | null, response: comm_pb.BasicResponse) => void): grpc.ClientUnaryCall;
    ping(request: comm_pb.BasicRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: comm_pb.BasicResponse) => void): grpc.ClientUnaryCall;
    ping(request: comm_pb.BasicRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: comm_pb.BasicResponse) => void): grpc.ClientUnaryCall;
    stop(request: comm_pb.BasicRequest, callback: (error: grpc.ServiceError | null, response: comm_pb.BasicResponse) => void): grpc.ClientUnaryCall;
    stop(request: comm_pb.BasicRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: comm_pb.BasicResponse) => void): grpc.ClientUnaryCall;
    stop(request: comm_pb.BasicRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: comm_pb.BasicResponse) => void): grpc.ClientUnaryCall;
    getValue(request: comm_pb.GetValueRequest, callback: (error: grpc.ServiceError | null, response: comm_pb.GetValueResponse) => void): grpc.ClientUnaryCall;
    getValue(request: comm_pb.GetValueRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: comm_pb.GetValueResponse) => void): grpc.ClientUnaryCall;
    getValue(request: comm_pb.GetValueRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: comm_pb.GetValueResponse) => void): grpc.ClientUnaryCall;
    setValue(request: comm_pb.SetValueRequest, callback: (error: grpc.ServiceError | null, response: comm_pb.SetValueResponse) => void): grpc.ClientUnaryCall;
    setValue(request: comm_pb.SetValueRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: comm_pb.SetValueResponse) => void): grpc.ClientUnaryCall;
    setValue(request: comm_pb.SetValueRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: comm_pb.SetValueResponse) => void): grpc.ClientUnaryCall;
    strlnValue(request: comm_pb.StrlnValueRequest, callback: (error: grpc.ServiceError | null, response: comm_pb.StrlnValueResponse) => void): grpc.ClientUnaryCall;
    strlnValue(request: comm_pb.StrlnValueRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: comm_pb.StrlnValueResponse) => void): grpc.ClientUnaryCall;
    strlnValue(request: comm_pb.StrlnValueRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: comm_pb.StrlnValueResponse) => void): grpc.ClientUnaryCall;
    deleteValue(request: comm_pb.DeleteValueRequest, callback: (error: grpc.ServiceError | null, response: comm_pb.DeleteValueResponse) => void): grpc.ClientUnaryCall;
    deleteValue(request: comm_pb.DeleteValueRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: comm_pb.DeleteValueResponse) => void): grpc.ClientUnaryCall;
    deleteValue(request: comm_pb.DeleteValueRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: comm_pb.DeleteValueResponse) => void): grpc.ClientUnaryCall;
    appendValue(request: comm_pb.AppendValueRequest, callback: (error: grpc.ServiceError | null, response: comm_pb.AppendValueResponse) => void): grpc.ClientUnaryCall;
    appendValue(request: comm_pb.AppendValueRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: comm_pb.AppendValueResponse) => void): grpc.ClientUnaryCall;
    appendValue(request: comm_pb.AppendValueRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: comm_pb.AppendValueResponse) => void): grpc.ClientUnaryCall;
    changeMembership(request: comm_pb.ChangeMembershipRequest, callback: (error: grpc.ServiceError | null, response: comm_pb.ChangeMembershipResponse) => void): grpc.ClientUnaryCall;
    changeMembership(request: comm_pb.ChangeMembershipRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: comm_pb.ChangeMembershipResponse) => void): grpc.ClientUnaryCall;
    changeMembership(request: comm_pb.ChangeMembershipRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: comm_pb.ChangeMembershipResponse) => void): grpc.ClientUnaryCall;
    appendEntries(request: comm_pb.AppendEntriesRequest, callback: (error: grpc.ServiceError | null, response: comm_pb.AppendEntriesResponse) => void): grpc.ClientUnaryCall;
    appendEntries(request: comm_pb.AppendEntriesRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: comm_pb.AppendEntriesResponse) => void): grpc.ClientUnaryCall;
    appendEntries(request: comm_pb.AppendEntriesRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: comm_pb.AppendEntriesResponse) => void): grpc.ClientUnaryCall;
    requestVote(request: comm_pb.RequestVoteRequest, callback: (error: grpc.ServiceError | null, response: comm_pb.RequestVoteResponse) => void): grpc.ClientUnaryCall;
    requestVote(request: comm_pb.RequestVoteRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: comm_pb.RequestVoteResponse) => void): grpc.ClientUnaryCall;
    requestVote(request: comm_pb.RequestVoteRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: comm_pb.RequestVoteResponse) => void): grpc.ClientUnaryCall;
}

export class CommServiceClient extends grpc.Client implements ICommServiceClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
    public ping(request: comm_pb.BasicRequest, callback: (error: grpc.ServiceError | null, response: comm_pb.BasicResponse) => void): grpc.ClientUnaryCall;
    public ping(request: comm_pb.BasicRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: comm_pb.BasicResponse) => void): grpc.ClientUnaryCall;
    public ping(request: comm_pb.BasicRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: comm_pb.BasicResponse) => void): grpc.ClientUnaryCall;
    public stop(request: comm_pb.BasicRequest, callback: (error: grpc.ServiceError | null, response: comm_pb.BasicResponse) => void): grpc.ClientUnaryCall;
    public stop(request: comm_pb.BasicRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: comm_pb.BasicResponse) => void): grpc.ClientUnaryCall;
    public stop(request: comm_pb.BasicRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: comm_pb.BasicResponse) => void): grpc.ClientUnaryCall;
    public getValue(request: comm_pb.GetValueRequest, callback: (error: grpc.ServiceError | null, response: comm_pb.GetValueResponse) => void): grpc.ClientUnaryCall;
    public getValue(request: comm_pb.GetValueRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: comm_pb.GetValueResponse) => void): grpc.ClientUnaryCall;
    public getValue(request: comm_pb.GetValueRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: comm_pb.GetValueResponse) => void): grpc.ClientUnaryCall;
    public setValue(request: comm_pb.SetValueRequest, callback: (error: grpc.ServiceError | null, response: comm_pb.SetValueResponse) => void): grpc.ClientUnaryCall;
    public setValue(request: comm_pb.SetValueRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: comm_pb.SetValueResponse) => void): grpc.ClientUnaryCall;
    public setValue(request: comm_pb.SetValueRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: comm_pb.SetValueResponse) => void): grpc.ClientUnaryCall;
    public strlnValue(request: comm_pb.StrlnValueRequest, callback: (error: grpc.ServiceError | null, response: comm_pb.StrlnValueResponse) => void): grpc.ClientUnaryCall;
    public strlnValue(request: comm_pb.StrlnValueRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: comm_pb.StrlnValueResponse) => void): grpc.ClientUnaryCall;
    public strlnValue(request: comm_pb.StrlnValueRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: comm_pb.StrlnValueResponse) => void): grpc.ClientUnaryCall;
    public deleteValue(request: comm_pb.DeleteValueRequest, callback: (error: grpc.ServiceError | null, response: comm_pb.DeleteValueResponse) => void): grpc.ClientUnaryCall;
    public deleteValue(request: comm_pb.DeleteValueRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: comm_pb.DeleteValueResponse) => void): grpc.ClientUnaryCall;
    public deleteValue(request: comm_pb.DeleteValueRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: comm_pb.DeleteValueResponse) => void): grpc.ClientUnaryCall;
    public appendValue(request: comm_pb.AppendValueRequest, callback: (error: grpc.ServiceError | null, response: comm_pb.AppendValueResponse) => void): grpc.ClientUnaryCall;
    public appendValue(request: comm_pb.AppendValueRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: comm_pb.AppendValueResponse) => void): grpc.ClientUnaryCall;
    public appendValue(request: comm_pb.AppendValueRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: comm_pb.AppendValueResponse) => void): grpc.ClientUnaryCall;
    public changeMembership(request: comm_pb.ChangeMembershipRequest, callback: (error: grpc.ServiceError | null, response: comm_pb.ChangeMembershipResponse) => void): grpc.ClientUnaryCall;
    public changeMembership(request: comm_pb.ChangeMembershipRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: comm_pb.ChangeMembershipResponse) => void): grpc.ClientUnaryCall;
    public changeMembership(request: comm_pb.ChangeMembershipRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: comm_pb.ChangeMembershipResponse) => void): grpc.ClientUnaryCall;
    public appendEntries(request: comm_pb.AppendEntriesRequest, callback: (error: grpc.ServiceError | null, response: comm_pb.AppendEntriesResponse) => void): grpc.ClientUnaryCall;
    public appendEntries(request: comm_pb.AppendEntriesRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: comm_pb.AppendEntriesResponse) => void): grpc.ClientUnaryCall;
    public appendEntries(request: comm_pb.AppendEntriesRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: comm_pb.AppendEntriesResponse) => void): grpc.ClientUnaryCall;
    public requestVote(request: comm_pb.RequestVoteRequest, callback: (error: grpc.ServiceError | null, response: comm_pb.RequestVoteResponse) => void): grpc.ClientUnaryCall;
    public requestVote(request: comm_pb.RequestVoteRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: comm_pb.RequestVoteResponse) => void): grpc.ClientUnaryCall;
    public requestVote(request: comm_pb.RequestVoteRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: comm_pb.RequestVoteResponse) => void): grpc.ClientUnaryCall;
}
