<script lang="ts">
    import { Link } from "svelte-routing";
</script>

<nav>
    <Link to="/">Summio</Link>
    <Link to="/docs">Docs</Link>
</nav>

<style>
    nav {
        display: flex;
        justify-content: left;
        align-items: center;
        padding: 1rem 1rem;
        background-color: #f5f5f5;
        border-bottom: 1px solid #e5e5e5;
        width:100%;
    }

    nav :global(a) {
        text-decoration: none;
        color: #333;
        margin-left: 1rem;
    }

    nav :global(a:hover) {
        text-decoration: underline;
    }
</style>