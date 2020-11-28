<script lang="ts">
    import type { Metadata } from "grpc-web";
    import { EnrollmentClient } from "../protos/JssoServiceClientPb";
    import { FinishEnrollmentRequest } from "../protos/jsso_pb";
    import { credentialFromJS } from "../lib/webauthn";
    import GrpcError from "../components/GrpcError.svelte";

    export let opts: PublicKeyCredentialCreationOptions;
    export let token: string;
    export let name: string;

    const enrollmentClient = new EnrollmentClient("", null, null);
    async function create() {
        const credential = await navigator.credentials.create({ publicKey: opts });
        if (!(credential instanceof PublicKeyCredential)) {
            throw "not a public key credential";
        }
        const req = new FinishEnrollmentRequest()
            .setCredential(credentialFromJS(credential))
            .setName(name);
        const metadata: Metadata = {};
        if (token != "") {
            metadata.authorization = "SessionID " + token;
        }
        return await enrollmentClient.finish(req, metadata);
    }
</script>

<style>
</style>

{#await create() then finishReply}
    <p>Successfully added a credential.</p>
    <p>
        This enrollment link is no longer valid, but you can
        <a href={finishReply.getLoginUrl()}>proceed to the login page</a>
        and login normally to add another authentication method.
    </p>
{:catch error}
    <p>There was a problem creating credential.</p>
    <GrpcError {error} />
{/await}
