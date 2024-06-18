<script>
  import { createEventDispatcher } from "svelte";
  import { slide, fade } from "svelte/transition";
  import { userInfo, allUsers } from "../../stores";
  import Button from "../../shared/button.svelte";
  import ImageToPost from "../../shared/imagePreview.svelte";
  const dispatch = createEventDispatcher();
  function closeOverlay() {
    dispatch("close");
  }
  console.log($userInfo);

  function autoResize() {
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

  let privatePost;
  let chooseUsers;

  let selectedUserIds = [];

  function togglePrivacy() {
    privatePost = !privatePost;
    chooseUsers = false;
    selectedUserIds = [];

    console.log(privatePost, "privatePost");
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
          <img src={$userInfo.avatar} alt="user avatar" />
          <p class="username">{$userInfo.firstName} {$userInfo.lastName}</p>
        </div>

        {#if chooseUsers}
          <select multiple bind:value={selectedUserIds}>
            {#each $allUsers as user}
              <option value={user.ID}>{user.FirstName} {user.LastName}</option>
            {/each}
          </select>
        {/if}

        <div class="privacy">
          {#if privatePost}
            <Button inverse={true} on:click={() => togglePrivacy()}
              >Set Public</Button
            >
            <Button
              type="secondary"
              inverse={true}
              on:click={() => toggleUsersList()}>Choose Users</Button
            >
          {:else}
            <Button
              type="secondary"
              inverse={true}
              on:click={() => togglePrivacy()}>Set Private</Button
            >
          {/if}
        </div>
      </div>
      <textarea on:input={autoResize} placeholder="What's on your mind?"
      ></textarea>
      <ImageToPost inputIDProp="postImage" fakeInputText="Add Image" />
      <div class="postButtons">
        <Button type="secondary">Post</Button>
        <button on:click={closeOverlay}>Cancel</button>
      </div>
    </div>
  </div>
</div>

<style>
  textarea {
    min-height: 200px;
  }

  select {
    margin-bottom: 0;
    border-color: greenyellow;
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
