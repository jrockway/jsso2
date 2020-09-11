<script lang="ts">
    export let name: string;
    let reply: Promise;
    export let sendName;
    function send() {
        reply = sendName(name).then((reply) => (name = reply.getResult()));
    }
</script>

<style>
    main {
        text-align: center;
        padding: 1em;
        max-width: 240px;
        margin: 0 auto;
    }

    h1 {
        color: #00aa22;
        text-transform: uppercase;
        font-size: 4em;
        font-weight: 100;
    }

    @media (min-width: 640px) {
        main {
            max-width: none;
        }
    }
</style>

<main>
    <h1>Hello {name}!</h1>
    <p>Type something: <input bind:value={name} /></p>
    <button on:click={send}>Submit</button>
    {#await reply}
        <p>Processing.</p>
    {:catch error}
        <p>Whoa! An error: {error.message}</p>
    {/await}
</main>
