jsso2 is a passwordless single-sign-on system written in Go and Typescript.

# Development

You can run everything locally. Run `npm run all`, and a server will begin listening on port 4000.
(Ignore the output from the various components that start themselves; those ports get proxied to by
the server running on port 4000.) Any Typescript/Svelte changes are reflected immediately as you
edit the files. Go changes require a restart for the time being. You will also need an `envoy`
binary. I get mine by `docker cp`-ing the binary out of one of
[their releases](https://www.envoyproxy.io/).

## Cleaning up Windows Hello keys

Windows remembers every enrollment you've ever done, which can result in confusion when you're
trying to log in. You can see a list of your Windows Hello keys by running `certutil -csp NGC -key`.
You should see keys that contain `FIDO_AUTHENTICATOR` in their names, and these are keys created by
Windows Hello. The part after `FIDO_AUTHENTICATOR` is the Relying Party ID hash and the User ID
separated by an underscore. Keys from the dev server will have a Relying Party ID of "localhost",
and the sha256sum of "localhost" is
49960de5880e8c687434170f6476605b8fe4aeb9a28632c7995cf3ba831d9763; this should allow you to recognize
the keys that you've enrolled while testing. You can then delete them with
`certutil -csp NGC -delkey <full key>` from an Administrator console.
