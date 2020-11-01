<script>
    import Router from "svelte-spa-router";
    import Header from "./components/Header.svelte";
    import Index from "./routes/Index.svelte";
    import Enrollment from "./routes/Enrollment.svelte";
    import NotFound from "./routes/NotFound.svelte";
    const routes = {
        "/": Index,
        "/enroll": Enrollment,
        "/enroll/:token": Enrollment,
        "*": NotFound,
    };
</script>

<body>
    <Header />
    <Router {routes} />
</body>
