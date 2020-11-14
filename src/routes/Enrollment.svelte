<script lang="ts">
    import type { Metadata } from "grpc-web";
    import { EnrollmentClient } from "../protos/JssoServiceClientPb";
    import { StartEnrollmentRequest } from "../protos/jsso_pb";
    import AddCredential from "../components/AddCredential.svelte";
    import { creationOptionsFromProto } from "../lib/webauthn";
    import GrpcError from "../components/GrpcError.svelte";

    const enrollmentClient = new EnrollmentClient("", null, null);

    export let params = {
        token: "",
    };
    let clicked = false;

    const metadata: Metadata = {};
    if (params.token != "") {
        metadata.authorization = "SessionID " + params.token;
    }

    let getUser = enrollmentClient
        .start(new StartEnrollmentRequest(), metadata)
        .then((response) => {
            if (response == null || !response.hasUser()) {
                throw "server error: no user in response";
            }
            if (response == null || !response.hasCredentialCreationOptions()) {
                throw "server error: no credential creation options in response";
            }
            return {
                username: response.getUser().getUsername(),
                opts: creationOptionsFromProto(response.getCredentialCreationOptions()),
            };
        });
</script>

<style>
</style>

<main>
    <h1>Enroll</h1>
    {#await getUser}
        <p>Validating your token.</p>
    {:then reply}
        <p>Welcome, {reply.username}!</p>
        <p>
            When you click the button below, your OS or browser will ask you to enroll a WebAuthn
            credential. We can't pick which one will be selected, but if you don't see the one you
            want pop up, pressing cancel will move on to the next one. When you find the one you
            want, enroll it. When you log in, you won't have to do this, and if you visit this
            enrollment page again, the credential you successfully enrolled will not be presented
            again.
        </p>
        <button on:click={() => (clicked = true)} disabled={clicked}>Enroll</button>
        {#if clicked}
            <AddCredential token={params.token} opts={reply.opts} />
        {/if}
    {:catch error}
        <p>There was a problem validating your token:</p>
        <GrpcError {error} />
    {/await}
</main>
