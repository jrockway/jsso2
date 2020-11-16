<script lang="ts">
    import { LoginClient } from "../protos/JssoServiceClientPb";
    import { StartLoginRequest } from "../protos/jsso_pb";
    import GrpcError from "../components/GrpcError.svelte";

    let username = "";
    let showLogin = true;
    let doStartLogin: Promise<StartLoginRequest>;

    const loginClient = new LoginClient("", null, null);

    function handleKeypress(event: KeyboardEvent) {
        if (event.key == "Enter") {
            startLogin();
        }
    }

    function startLogin() {
        const req = new StartLoginRequest();
        doStartLogin = loginClient.start(req, null);
        showLogin = false;
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
                on:keydown={handleKeypress} /><button on:click={startLogin}>Login</button>
        </p>
    {:else}
        {#await doStartLogin}
            Hello, <b>{username}</b>.
        {:then reply}
            Here's the reply: {reply}.
        {:catch error}
            <p>There was a problem beginning the login process.</p>
            <GrpcError {error} />
        {/await}
    {/if}
</main>
