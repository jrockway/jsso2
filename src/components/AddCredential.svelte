<script lang="ts">
    const randomStringFromServer = "foobar";
    let opts: PublicKeyCredentialCreationOptions = {
        challenge: Uint8Array.from(randomStringFromServer, (c) => c.charCodeAt(0)),
        rp: {
            name: "Example",
            id: "localhost",
        },
        user: {
            id: Uint8Array.from("foobar", (c) => c.charCodeAt(0)),
            name: "username",
            displayName: "Username",
        },
        pubKeyCredParams: [
            {
                type: "public-key",
                alg: -7,
            },
            {
                type: "public-key",
                alg: -35,
            },
            {
                type: "public-key",
                alg: -36,
            },
            {
                type: "public-key",
                alg: -257,
            },
            {
                type: "public-key",
                alg: -258,
            },
            {
                type: "public-key",
                alg: -259,
            },
            {
                type: "public-key",
                alg: -37,
            },
            {
                type: "public-key",
                alg: -38,
            },
            {
                type: "public-key",
                alg: -39,
            },
            {
                type: "public-key",
                alg: -8,
            },
        ],
        authenticatorSelection: {
            //authenticatorAttachment: "cross-platform",
            userVerification: "discouraged",
        },
        timeout: 60000,
        attestation: "direct",
    };

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
