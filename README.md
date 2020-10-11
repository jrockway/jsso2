jsso2 is a passwordless single-sign-on system written in Go and Typescript.

# Development

You can run everything locally. Run `npm run all`, and a server will begin listening on port 4000.
(Ignore the output from the various components that start themselves; those ports get proxied to by
the server running on port 4000.) Any Typescript/Svelte changes are reflected immediately as you
edit the files. Go changes require a restart for the time being. You will also need an `envoy`
binary. I get mine by `docker cp`-ing the binary out of one of
[their releases](https://www.envoyproxy.io/).
