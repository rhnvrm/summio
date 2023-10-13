<script lang="ts">
    import {success, warning, failure} from "./toast"
    import pulse from "../assets/pulse.svg";

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
                if (result.status == "success") {
                    success("File uploaded successfully")
                    // wair for 2 seconds then redirect to the new doc
                    setTimeout(() => {
                        window.location.href = "/docs/" + result.data.id
                    }, 2000)
                } else {
                    failure("Error uploading file: " + result.message)
                }
            })
            .catch((error) => {
                console.error("Error:", error);
                failure("Error uploading file: " + error)
            });
        inProgress = false;
    }
</script>

<div class="uploader box">
    <h1>Upload</h1>
    <div class="form">
    <input name="file" id="file" type="file" bind:files />

    <button disabled={inProgress} on:click={upload}>Submit</button>

    {#if inProgress}
    <img alt="loading" src={pulse} />   
    {/if}
    </div>
</div>

<style>
    .form {
        display: flex;
        flex-direction: row;
    }
</style>
