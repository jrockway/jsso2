<script lang="ts">
    export let opts: PublicKeyCredentialCreationOptions;

    let result = navigator.credentials.create({ publicKey: opts }).then((credential) => {
        if (credential instanceof PublicKeyCredential) {
            return credential;
        }
        throw "not a public key credential";
    });
</script>

<style>
</style>

{#await result then credential}
    <p>Here is your credential:</p>
    <table>
        <thead>
            <tr>
                <td>Id</td>
                <td>Type</td>
                <td>Response</td>
            </tr>
        </thead>
        <tr>
            <td>{credential.id}</td>
            <td>{credential.type}</td>
            <td>{credential.response}</td>
        </tr>
    </table>
{:catch error}
    <p>Error: {error}</p>
{/await}
