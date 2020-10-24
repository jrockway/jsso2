<script>
    import { location } from "svelte-spa-router";
</script>

<style>
</style>

<header>
    <a href="/#/">Home</a>
    {#if $location.length > 0}&gt; {$location.slice(1)}{/if}
</header>
