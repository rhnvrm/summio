<script lang="ts">
    import type { PDFDocType } from "../types/docs.type.ts";

    import pulse from "../assets/pulse.svg";
    let docs: PDFDocType[] = [];

    let searchQuery = "";
    let inProgress = false;

    async function search() {
        inProgress = true;
        await fetch("/api/pdf?q=" + searchQuery + "&limit=10")
            .then((res) => res.json())
            .then((res) => {
                docs = res.data;
            });
        inProgress = false;
    }
</script>

<div class="list_docs box">
    <div class="form">
        <input type="text" bind:value={searchQuery} />

        <button disabled={inProgress} on:click={search}>Search</button>
    </div>

    {#if inProgress}
        <img alt="loading" src={pulse} />
    {/if}
    <h1>Search Results</h1>
    {#if docs === null || docs.length === 0}
        <p>No docs found</p>
    {:else}
        <ul>
            {#each docs as doc}
                <li>
                    <a href="/docs/{doc.id}">{doc.title}</a>
                </li>
            {/each}
        </ul>
    {/if}
</div>

<style>
    .form {
        display: flex;
        flex-direction: row;
    }
</style>
