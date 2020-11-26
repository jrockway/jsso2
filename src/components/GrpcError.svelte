<script lang="ts">
    export let error: any;
</script>

<style>
</style>

{#if error.message !== undefined && error.code !== undefined}
    <p class="error">{error.message} ({error.code})</p>
{:else}
    <p class="error">{error}</p>
{/if}
