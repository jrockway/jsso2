<script lang="ts">
    import type { Metadata } from "grpc-web";
    import { EnrollmentClient } from "../protos/JssoServiceClientPb";
    import { StartEnrollmentRequest } from "../protos/jsso_pb";
    import AddCredential from "../components/AddCredential.svelte";

    const enrollmentClient = new EnrollmentClient("", null, null);

    export let params = {
        token: "",
    };

    const metadata: Metadata = {};
    if (params.token) {
        metadata.authorization = "SessionID " + params.token;
    }

    let getUser = enrollmentClient
        .start(new StartEnrollmentRequest(), metadata)
        .then((response) => {
            if (response == null || response.getUser() == null) {
                throw "no user in response";
            }
            return response;
        });
</script>

<style>
</style>

<main>
    <h1>Enroll</h1>
    {#await getUser}
        <p>Validating your token.</p>
    {:then reply}
        <p>Welcome, {reply.getUser().getUsername()}!</p>
        <AddCredential />
    {:catch error}
        <p>We can't validate your token:</p>
        {#if error.message != null}
            <p>{error.message} ({error.code})</p>
        {:else}
            <p>{error}</p>
        {/if}
    {/await}
</main>
