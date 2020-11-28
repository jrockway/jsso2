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

    async function getUser() {
        const reply = await enrollmentClient.start(new StartEnrollmentRequest(), metadata);
        if (reply == null || !reply.hasUser()) {
            throw "server error: no user in response";
        }
        if (reply == null || !reply.hasCredentialCreationOptions()) {
            throw "server error: no credential creation options in response";
        }
        return {
            username: reply.getUser().getUsername(),
            opts: creationOptionsFromProto(reply.getCredentialCreationOptions()),
        };
    }
</script>

<style>
</style>

<main>
    <h1>Enroll</h1>
    {#await getUser()}
        <p>Validating your token.</p>
    {:then reply}
        <p>Welcome, <b>{reply.username}</b>!</p>
        <p>
            When you click the button below, your OS or browser will ask you to enroll a WebAuthn
            credential. We can't pick which one will be selected, but if you don't see the one you
            want pop up, pressing cancel will move on to the next one. When you find the one you
            want, enroll it. When you log in, you won't have to do this.
        </p>
        <p>
            Note that for privacy reasons, even if you have already enrolled a key, your OS or
            browser will pretend that you've never enrolled that authenticator before, let you
            enroll it, and then give you an error ("The user attempted to register an authenticator
            that contains one of the credentials already registered with the relying party.").
            There's nothing we can do about that. We tell your browser which credentials we already
            have, but it ignores them until you do the auth dance.
        </p>
        {#if !clicked}
            <button id="enroll" on:click={() => (clicked = true)}>Enroll</button>
        {:else}
            <AddCredential token={params.token} opts={reply.opts} />
        {/if}
    {:catch error}
        <p>There was a problem validating your token.</p>
        <GrpcError {error} />
    {/await}
</main>
