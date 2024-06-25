<script>
  import { selectUser } from "../../utils";
  import { onMount } from "svelte";
  export let comment;
</script>

<div class="singleComment">
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-mouse-events-have-key-events -->
  <div
    on:mouseover={(event) => {
      let createdAt = event.currentTarget.querySelector(".commentCreatedAt");
      createdAt.style.display = "inline";
    }}
    on:mouseout={(event) => {
      let createdAt = event.currentTarget.querySelector(".commentCreatedAt");
      createdAt.style.display = "none";
    }}
    class="userInfo"
    on:click={() => selectUser(comment.userID)}
  >
    <p class="commentCreator">
      {comment.user.firstName}
      {comment.user.lastName}
    </p>
    <p class="commentCreatorAvatar">
      <img
        class="postCreatorAvatar"
        src={comment.user.avatar}
        alt="user avatar"
      />
    </p>
    <p class="commentCreatedAt">
      {comment.createdAt}
    </p>
  </div>

  <div class="commentContent">
    {comment.content}
    {#if comment.img}
      <img src={comment.img} alt="This is a comment" />
    {/if}
  </div>
</div>

<style>

div {
    padding: 8px;
    border-radius: 8px;
  }

  p {
    margin: 0;
  }
  
  .postCreatorAvatar {
    padding: 4px 0;
    max-height: 40px;
  }
  .userInfo {
    grid-area: userInfo;
    cursor: pointer;
  }
  .singleComment {
    font-size: small;
    border: solid 1px #333;
    display: grid;
    grid-auto-columns: 1fr;
    grid-template-columns: 0.3fr 0.3fr 1.5fr;
    grid-template-rows: 0.5fr 0.5fr 3fr;
    grid-template-areas:
      ". userInfo commentContent"
      ". userInfo commentContent"
      ". userInfo commentContent";
  }

  .singleComment img {
    padding: 12px;
  }

  .commentCreatedAt {
    display: none;
  }
  .commentContent {
    display: flex;
    flex-direction: column;
    grid-area: commentContent;
    margin: 0 8px;
    border: solid 1px #333;
  }
  .commentContent > img {
    max-width: 300px;
    max-height: 300px;
  }
</style>
