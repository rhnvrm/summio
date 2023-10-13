<script lang="ts">
    import ViewDoc from "../lib/view_doc.svelte";

    export let docID: string;
</script>

<div>
    <ViewDoc docID="{docID}" />
</div>
