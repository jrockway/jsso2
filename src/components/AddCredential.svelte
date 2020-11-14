<script lang="ts">
    import type { Metadata } from "grpc-web";
    import { EnrollmentClient } from "../protos/JssoServiceClientPb";
    import { FinishEnrollmentRequest } from "../protos/jsso_pb";
    import { credentialFromJS } from "../lib/webauthn";
    import GrpcError from "../components/GrpcError.svelte";

    export let opts: PublicKeyCredentialCreationOptions;
    export let token: string;

    const enrollmentClient = new EnrollmentClient("", null, null);
    let creationResult = navigator.credentials.create({ publicKey: opts }).then((credential) => {
        if (credential instanceof PublicKeyCredential) {
            return credential;
        }
        throw "not a public key credential";
    });
    let submissionResult = creationResult.then((credential) => {
        const req = new FinishEnrollmentRequest().setCredential(credentialFromJS(credential));
        const metadata: Metadata = {};
        if (token != "") {
            metadata.authorization = "SessionID " + token;
        }
        return enrollmentClient.finish(req, metadata);
    });
</script>

<style>
</style>

{#await creationResult then credential}
    {#await submissionResult}
        <p>Your credential "{credential.id}" has been created; sending it to the server...</p>
    {:then}
        <p>Successfully submitted the credential.</p>
    {:catch error}
        <p>There was a problem sending the credential to the server:</p>
        <GrpcError {error} />
    {/await}
{:catch error}
    <p>There was a problem creating the credential:</p>
    <p>{error}</p>
{/await}
