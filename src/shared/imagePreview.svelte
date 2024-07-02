<script>
  import { displayUserAuthError, uploadImageStore, userInfo } from "../stores";
  import Button from "./button.svelte";

  //PROPS
  export let fakeInputText = "";
  export let fakeInputMaxAvatarSize = "[Max: 500KB]";
  export let inputIDProp = "";
  export let style = "";
  export let src;

  let input;
  let deleteImg = false;
  let image;
  let showImage = false;
  if (src) {
    showImage = true;
  }
  const maxFileSize = 500 * 1024; // 500 KB
  const allowedTypes = ["image/jpeg", "image/jpg", "image/png", "image/gif"];

  uploadImageStore.set(uploadImage);

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
      deleteImg = false;
      const reader = new FileReader();
      reader.addEventListener("load", function () {
        image.setAttribute("src", reader.result);
        image.setAttribute("name", file.name);
      });
      reader.readAsDataURL(file);
      // Shorten filename if needed
      let filename = file.name;
      let fileExtension = filename.split(".").pop();
      if (filename.length > 17) {
        filename = filename.slice(0, 17) + "...." + fileExtension;
      }

      fakeInputText = filename + " (Change)";
      fakeInputMaxAvatarSize = (file.size / 1024).toFixed(0) + "KB";
      return;
    }
    deleteImg = true;
    showImage = false;
  }

  export function removeImage() {
    showImage = false;
    image.src = "";
    image.name = "";
    fakeInputText = "Try another image";
    fakeInputMaxAvatarSize = "[Max: 500KB]";
    input.value = null;
    deleteImg = true;
  }

  // anti upload image :((
  export async function uploadImage(obj) {
    const file = input.files[0];
    if (file) {
      console.log(file);
      const formData = new FormData();
      formData.append("image", file);
      formData.append("from", inputIDProp); // From which prop id the upload is coming from
      if (obj) {
        formData.append("postID", obj.post); // Should be generated somehow with the new post ID
        formData.append("commentID", obj.comment); // Should be generated somehow with the comment iD
      }
      const response = await fetch("/uploadImage", {
        method: "POST",
        body: formData,
      });
      let path = await response.text();
      if (obj) removeImage();
      return path;
    }
    if (deleteImg) {
      return "";
    }
  }
</script>

<label class="fakeInput" {style} for={inputIDProp}
  >{fakeInputText}
  <p class="maxImageSize">{fakeInputMaxAvatarSize}</p>
</label>
<input
  id={inputIDProp}
  bind:this={input}
  on:change={onChange}
  type="file"
  style="display:none"
/>

{#if showImage}
  <div>
    <Button inverse={true} on:click={removeImage} customStyle="width:100%;"
      >Remove image</Button
    >
    <img id="imagePreview" bind:this={image} {src} alt="Preview" />
  </div>
{/if}

<style>
  .maxImageSize {
    width: max-content;
    margin: 0;
    margin-left: 10px;
    font-size: 13px;
    color: #777;
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
    max-width: 300px;
  }

  div {
    max-width: 300px;
    min-height: 100px;
    margin: 8px auto;
  }
  img {
    width: 100%;
  }
</style>
