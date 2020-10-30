<script>
    import { location } from "svelte-spa-router";
</script>

<style>
</style>

<header>
    <a href="/#/">Home</a>
    {#if $location.length > 1}&gt; {$location.split('/')[1]}{/if}
</header>
