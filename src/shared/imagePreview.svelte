<script>
    
    let input;
    let image;
    let showImage = false;
  
    function onChange() {
        const file = input.files[0];
            
        if (file) {
            showImage = true;
            const reader = new FileReader();
            reader.addEventListener("load", function () {
                image.setAttribute("src", reader.result);
            });
            reader.readAsDataURL(file);

            // Shorten filename if needed
            let filename = file.name;
            let fileExtension = filename.split('.').pop();
            if (filename.length > 17) {
                filename = filename.slice(0, 17) + '....' + fileExtension;
            }
            const container = document.querySelector('label[for="avatar"]');
            container.textContent = filename + ' (Change)';
                    
            return;
        }
        showImage = false; 

        const container = document.querySelector('label[for="avatar"]');
        container.textContent = 'Upload avatar (Optional)';
    }
    
</script>

<label class="fakeInput" for="avatar">Upload avatar (Optional)</label>
<input id="avatar"
    bind:this={input}
    on:change={onChange}
    type="file"
    style="display:none"
/>
<!-- bind:value={userData.avatar} -->
{#if showImage}
    <div>
        <img id="avatarPreview" bind:this={image} src="" alt="Preview" />
    </div>
{/if}

<style>

    .fakeInput {
    color: #ddd;
    display: inline-block;
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
