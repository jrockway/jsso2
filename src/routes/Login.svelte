<script lang="ts">
    import { LoginClient } from "../protos/JssoServiceClientPb";
    import { StartLoginRequest, StartLoginReply } from "../protos/jsso_pb";
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
        console.log(reply);
        const details = await navigator.credentials.get({});
        return details.type;
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
        {:then reply}
            Here's the reply: {reply}.
        {:catch error}
            <p>There was a problem beginning the login process.</p>
            <GrpcError {error} />
        {/await}
    {/if}
</main>
