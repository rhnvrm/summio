<script lang="ts">
    import { onMount } from "svelte";
    import SvelteMarkdown from "svelte-markdown";
    import type { PDFDocType } from "../types/docs.type";

    export let docID: string;

    let doc: PDFDocType;

    onMount(() => {
        fetch("/api/pdf/" + docID)
            .then((res) => res.json())
            .then((res) => {
                console.log(res);
                doc = res.data;
            });
    });
</script>

<div class="box">
    {#if doc}
        <h1>{doc.title}</h1>

        <div class="summary">
            <span>AI Generated Summary</span>
            <SvelteMarkdown source={doc.summary} />
        </div>

        <iframe
            title="doc"
            src="/api/static/docs/{doc.file}"
            width="100%"
            height="800px"
        />

        {#if doc.intermediate_summary && doc.intermediate_summary.length > 0}
            <details>
                <h2>Details</h2>
                <ul>
                    {#each doc.intermediate_summary as summary}
                        <li>
                            {summary}
                        </li>
                    {/each}
                </ul>
            </details>
        {/if}
    {/if}
    <!-- <h1>{doc.title}</h1> -->
</div>

<style>
    .summary span {
        font-weight: light;
        color:rgb(93, 93, 93);
    }
</style>