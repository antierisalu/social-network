<script>
  import { createEventDispatcher } from "svelte";
  import { slide, fade } from "svelte/transition";
  import {
    userInfo,
    allUsers,
    uploadImageStore,
    groupSelected,
    API_URL,
    IMAGE_URL,
  } from "../../stores";
  import Button from "../../shared/button.svelte";
  import ImageToPost from "../../shared/imagePreview.svelte";
  import { getPosts } from "../../utils";

  const dispatch = createEventDispatcher();
  function closeOverlay() {
    dispatch("close");
  }

  function autoResize() {
    // for automatic resize of post content textarea
    const maxHeight = window.innerHeight * 0.8;
    const minHeight = 200;
    this.style.height = "auto";
    if (this.scrollHeight > maxHeight) {
      this.style.height = maxHeight + "px";
      this.style.overflowY = "scroll";
    } else if (this.scrollHeight < minHeight) {
      this.style.height = minHeight + "px";
      this.style.overflowY = "hidden";
    } else {
      this.style.height = this.scrollHeight + "px";
      this.style.overflowY = "hidden";
    }
  }

  let privatePost = false;
  let chooseUsers = false;
  let selectedUserIds;
  let content = "";

  let uploadImage;
  uploadImageStore.subscribe((value) => {
    uploadImage = value;
  });

  $: post = {
    userID: $userInfo.id,
    content: content,
    img: "",
    privacy: Number(privatePost + chooseUsers),
    groupID: $groupSelected,
    customPrivacyIDs: selectedUserIds,
  };

  async function sendPost() {
    console.log(post);
    if (!post.content) {
      alert("Post cannot be empty");
      return;
    }
    const response = await fetch(`${API_URL}/newPost`, {
      credentials: "include",
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        userID: post.userID,
        Content: post.content,
        Img: post.img,
        GroupID: $groupSelected,
        Privacy: post.privacy,
        CUSTOMPrivacyIDs: post.customPrivacyIDs,
      }),
    });

    const postID = await response.json();
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    closeOverlay();
    await uploadImage({ post: postID }).catch((error) => {
      console.error("Error uploading the image:", error);
    });
    getPosts();
  }

  function togglePrivacy() {
    privatePost = !privatePost;
    chooseUsers = false;
    selectedUserIds = [];

    console.log(privatePost, "privatePost");
    console.log($userInfo);
  }

  function toggleUsersList() {
    chooseUsers = !chooseUsers;
  }
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div out:fade={{ duration: 150 }} class="overlay"></div>
<div in:slide out:fade={{ duration: 150 }} class="modal">
  <div class="modal-content">
    <div class="createPost">
      <div class="overlayHeader">
        <div class="userInfo">
          <img src={IMAGE_URL}{$userInfo.avatar} alt="user avatar" />
          <p class="username">{$userInfo.firstName} {$userInfo.lastName}</p>
        </div>

        {#if chooseUsers}
          <div>
            Select Users
            <select multiple bind:value={selectedUserIds}>
              {#each $userInfo.followers as user}
                <option value={user.ID}>{user.FirstName} {user.LastName}</option
                >
              {/each}
            </select>
          </div>
        {/if}
        {#if $groupSelected === 0}
          <div class="privacy">
            {#if privatePost}
              <Button inverse={true} on:click={() => togglePrivacy()}
                >Set Public</Button
              >
              {#if $userInfo.followers}
                <Button
                  type="secondary"
                  inverse={true}
                  on:click={() => toggleUsersList()}
                >
                  {#if chooseUsers}Regular Privacy
                  {:else}Custom Privacy
                  {/if}</Button
                >
              {/if}
            {:else}
              <Button
                type="secondary"
                inverse={true}
                on:click={() => togglePrivacy()}>Set Private</Button
              >
            {/if}
          </div>
        {/if}
      </div>
      <textarea
        on:input={autoResize}
        bind:value={content}
        placeholder="What's on your mind?"
      ></textarea>
      <ImageToPost inputIDProp="postImage" fakeInputText="Add Image" />
      <div class="postButtons">
        <Button type="secondary" on:click={() => sendPost()}>Post</Button>
        <Button on:click={closeOverlay}>Cancel</Button>
      </div>
    </div>
  </div>
</div>

<style>
  textarea {
    min-height: 200px;
  }

  select {
    margin-top: 5px;
    margin-bottom: 5px;
    border-color: greenyellow;
    height: 100px;
    scrollbar-width: thin;
    scrollbar-color: greenyellow #011;
  }

  .privacy {
    display: grid;
    grid-row: auto;
  }

  .overlayHeader {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
    min-height: 100px;
  }

  .userInfo {
    display: grid;
    grid-template-columns: 50px auto;
    grid-template-rows: 30px 20px;
    padding: 5px;
    align-items: center;
  }

  p {
    text-align: left;
    margin-left: 50px;
  }

  .username {
    grid-row: 1;
    grid-column: 4;
  }

  .createPost {
    display: flex;
    flex-direction: column;
    border-radius: 16px;
  }
  img {
    margin-right: 10px;
    grid-column: 1;
    grid-row: 1/3;
    border-radius: 50px;
    max-width: 90px;
  }

  .overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    z-index: 0;
  }

  .modal {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background-color: #011;
    padding: 20px;
    box-shadow: 0 5px 15px rgba(0, 0, 0, 0.3);
    z-index: 2;
    width: 80%;
    max-width: 500px;
    border-radius: 8px;
  }

  .modal-content {
    position: relative;
  }
</style>
