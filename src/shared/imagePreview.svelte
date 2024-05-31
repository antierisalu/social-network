<script>
    
    let input;
    let container;
    let image;
    let placeholder;
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
                    
            return;
        }
        showImage = false; 
    }
</script>
        <label class="fakeInput" for="avatar">Upload avatar (Optional)</label>
        <input id="avatar"
            bind:this={input}
            on:change={onChange}
        type="file" style="display:none"
        />
        <!-- bind:value={userData.avatar} -->
        {#if showImage}
            <div bind:this={container}>
                <img id="avatarPreview" bind:this={image} src="" alt="Preview" />
            </div>
        {/if}

<!-- 
    <div class="fakeInput">
        <label for="avatar">Upload avatar (Optional)</label>
        <input type="file" id="avatar" class="hidden" accept="image/*" bind:value={userData.avatar}>
    </div> -->

<style>


    .fakeInput {
    color: #ddd;
    display: inline-block;
    padding: 8px 10px;
    border: 1px solid #ccc;
    border-radius: 6px;
    background-color: #011;
    margin-bottom: 8px;
    cursor: pointer;
    }

    .fakeInput label {
    cursor: pointer;
    text-align: left;
    margin: 1;
    opacity: 0.7; 
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
