<script lang="ts">
    let files: { name: any }[];
    let inProgress = false;

    async function upload() {
        inProgress = true;
        const formData = new FormData();
        formData.append("file", files[0]);
        await fetch("/api/pdf", {
            method: "POST",
            body: formData,
        })
            .then((response) => response.json())
            .then((result) => {
                console.log("Success:", result);
            })
            .catch((error) => {
                console.error("Error:", error);
            });
        inProgress = false;
    }
</script>

<div class="uploader box">
    <h1>Upload</h1>
    <input name="file" id="file" type="file" bind:files />

    <button disabled={inProgress} on:click={upload}>Submit</button>
</div>

<style>

</style>