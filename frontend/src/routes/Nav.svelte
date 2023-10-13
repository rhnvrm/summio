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
        color: #454545;
        background-color: #f5f5f5;
        border-bottom: 1px solid #e5e5e5;
        box-shadow: rgba(60, 64, 67, 0.3) 0px 1px 2px 0px,
            rgba(60, 64, 67, 0.15) 0px 1px 3px 1px;
        width: 100%;
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
