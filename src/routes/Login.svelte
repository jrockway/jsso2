<script lang="ts">
    import { LoginClient } from "../protos/JssoServiceClientPb";
    import { StartLoginRequest } from "../protos/jsso_pb";
    import { requestOptionsFromProto } from "../lib/webauthn";
    import GrpcError from "../components/GrpcError.svelte";

    let username = "";
    let showLogin = true;

    const loginClient = new LoginClient("", null, null);

    function handleKeypress(event: KeyboardEvent) {
        if (event.key == "Enter") {
            showLogin = false;
        }
    }

    async function login(u: string) {
        showLogin = false;
        const req = new StartLoginRequest();
        req.setUsername(u);
        const reply = await loginClient.start(req, null);
        console.log(reply.toObject);
        const publicKey = requestOptionsFromProto(reply.getCredentialRequestOptions());
        publicKey.userVerification = "discouraged";
        const assertion = await navigator.credentials.get({
            publicKey: publicKey,
        });
        console.log(assertion);
        return assertion;
    }
</script>

<style>
</style>

<main>
    <h1>Login</h1>
    {#if showLogin}
        <p>
            Enter your username: <input
                type="text"
                bind:value={username}
                on:keydown={handleKeypress} /><button
                on:click={() => {
                    showLogin = false;
                }}>Login</button>
        </p>
    {:else}
        {#await login(username)}
            Hello, <b>{username}</b>.
        {:then assertion}
            Here's the reply: {assertion}.
        {:catch error}
            <p>There was a problem logging in.</p>
            <GrpcError {error} />
        {/await}
    {/if}
</main>
