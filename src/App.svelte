<script>
    import Router from "svelte-spa-router";
    import Header from "./components/Header.svelte";
    import Index from "./routes/Index.svelte";
    import GrpcWebDemo from "./routes/GrpcWebDemo.svelte";
    import NotFound from "./routes/NotFound.svelte";
    const routes = {
        "/": Index,
        "/grpc-web-demo": GrpcWebDemo,
        "*": NotFound,
    };
</script>

<body>
    <Header />
    <Router {routes} />
</body>
