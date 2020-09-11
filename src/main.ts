import App from "./App.svelte";

import * as grpcWeb from "grpc-web";
import { NameServiceClient } from "./protos/FooServiceClientPb";
import { TransformNameRequest, TransformNameReply } from "./protos/foo_pb";
const nameServiceClient = new NameServiceClient("http://localhost:4000", null, null);

function sendName(name: string): Promise<TransformNameReply> {
    let req = new TransformNameRequest();
    req.setName(name);
    return nameServiceClient.transformName(req, null);
}

const app = new App({
    target: document.body,
    props: {
        name: "world",
        sendName: sendName,
    },
});

export default app;
