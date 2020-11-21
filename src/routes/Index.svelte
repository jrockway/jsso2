<script lang="ts">
    import { UserClient } from "../protos/JssoServiceClientPb";
    import { WhoAmIRequest } from "../protos/jsso_pb";
    import GrpcError from "../components/GrpcError.svelte";

    const userClient = new UserClient("", null, null);
    const currentUser = userClient.whoAmI(new WhoAmIRequest(), null).then((reply) => {
        return reply.getUser();
    });
</script>

<h1>Welcome</h1>
{#await currentUser}
    <p>Getting details about your session...</p>
{:then user}
    {#if user === undefined}
        <p>You're not logged in. Proceed to <a href="/#/login">the login page</a>.</p>
    {:else}
        <p>Welcome, <strong>{user.getUsername()}</strong>.</p>
        <p>You can manage your account from this page.</p>
        <ul>
            <li><a href="/#/enroll">Enroll a security key</a>.</li>
            <li><a href="/logout">Log out</a>.</li>
        </ul>
    {/if}
{:catch error}
    <p>There was a problem checking the status of your session. Reload this page to try again!</p>
    <GrpcError {error} />
{/await}
