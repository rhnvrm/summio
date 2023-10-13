<script lang="ts">
    import { onMount } from "svelte";
    import type { PDFDocType } from "../types/docs.type.ts";

    let docs: PDFDocType[] = [];

    onMount(() => {
        fetch("/api/pdf")
            .then((res) => res.json())
            .then((res) => {
                docs = res.data;
            });
    });
</script>

<div class="list_docs box">
    <h1>Docs</h1>
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
</style>
