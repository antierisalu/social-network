<script>
    import { displayUserAuthError } from "../stores";
    import Button from "./button.svelte";

    export let fakeInputText
    export let fakeInputMaxAvatarSize = '[Max: 500KB]'

    let input;
    let image;
    let showImage = false;
    const maxFileSize = 500 * 1024; // 500 KB
    const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif'];
  
    function onChange() {
        const file = input.files[0];
        if (file.size > maxFileSize) {
            displayUserAuthError("File size too big!");
            return;
        }
        if (!allowedTypes.includes(file.type)) {
            displayUserAuthError("Please use jpeg, jpg, png or gif");
            return;
        }
        if (file) {
            showImage = true;
            const reader = new FileReader();
            reader.addEventListener("load", function () {
                image.setAttribute("src", reader.result);
                image.setAttribute("name", file.name);
            });
            reader.readAsDataURL(file);
            // Shorten filename if needed
            let filename = file.name;
            let fileExtension = filename.split('.').pop();
            if (filename.length > 17) {
                filename = filename.slice(0, 17) + '....' + fileExtension;
            }
            const container = document.querySelector('label[for="uploadedImage"]');
            container.textContent = filename + ' (Change)';
            return;
        }
        showImage = false; 
    }

    function removeImage() {
        showImage = false;
        image.src = "";
        image.name = "";
        const container = document.querySelector('label[for="uploadedImage"]');
        container.textContent = fakeInputText + ' ' + fakeInputMaxAvatarSize
        input.value = null;
    }
    
</script>

<label class="fakeInput" for="uploadedImage">{fakeInputText}
    <p class="maxImageSize">{fakeInputMaxAvatarSize}</p>
</label>
<input id="uploadedImage"
    bind:this={input}
    on:change={onChange}
    type="file"
    style="display:none"
/>
{#if showImage}
    <div>
        <Button inverse={true} on:click={removeImage} customStyle="width:100%;">Remove image</Button>
        <img id="imagePreview" bind:this={image} src="" alt="Preview" />
    </div>
{/if}

<style>
    .maxImageSize {
        width: max-content;
        margin: 0;
        margin-left: 10px;
        font-size: 13px;
        color: #777
    }

    .fakeInput {
    color: #ddd;
    display: inline-flex;
    align-items: center;
    text-align: left;
    padding: 8px 10px;
    border: 1px solid #ccc;
    border-radius: 6px;
    background-color: #011;
    margin-bottom: 8px;
    cursor: pointer;
    }

    input {
        width: 300px;
    }

    div {
    width: 300px;
    min-height: 100px;
    margin: 8px auto;
 
    } 
     img {
        width: 100%;
    }
</style>
