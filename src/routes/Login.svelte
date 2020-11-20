<script lang="ts">
    import { LoginClient } from "../protos/JssoServiceClientPb";
    import { StartLoginRequest, FinishLoginRequest } from "../protos/jsso_pb";
    import { credentialFromJS, requestOptionsFromProto } from "../lib/webauthn";
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
        const startReq = new StartLoginRequest();
        startReq.setUsername(u);
        const startReply = await loginClient.start(startReq, null);
        const publicKey = requestOptionsFromProto(startReply.getCredentialRequestOptions());
        publicKey.userVerification = "discouraged";
        const assertion = await navigator.credentials.get({
            publicKey: publicKey,
        });
        if (!(assertion instanceof PublicKeyCredential)) {
            throw "not a public key credential";
        }
        const finishReq = new FinishLoginRequest();
        finishReq.setCredential(credentialFromJS(assertion));
        const finishReply = await loginClient.finish(finishReq, {
            Authorization: "SessionID " + startReply.getToken(),
        });
        return finishReply.getRedirectUrl();
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
            <p>Hello, <b>{username}</b>.</p>
        {:then redirect}
            <p>You have logged in.</p>
            {#if redirect != ''}
                <p>You should be redirected to <a href={redirect}>{redirect}</a> shortly.</p>
            {/if}
        {:catch error}
            <p>There was a problem logging in.</p>
            <GrpcError {error} />
        {/await}
    {/if}
</main>
