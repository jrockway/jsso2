<script lang="ts">
    export let error: any;
</script>

<style>
</style>

{#if error.message !== undefined && error.code !== undefined}
    <p>{error.message} ({error.code})</p>
{:else}
    <p>{error}</p>
{/if}
