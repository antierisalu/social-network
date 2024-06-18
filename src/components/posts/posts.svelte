<script>
  import Button from "../../shared/button.svelte";
  import ImageToPost from "../../shared/imagePreview.svelte";
  import { slide } from "svelte/transition";
  import PostOverlay from "./createPost.svelte";
  import ImageToComment from "../../shared/imagePreview.svelte";
  import { userInfo } from "../../stores";

  let showOverlay = false;
  let showComments = false;
  const openProfile = () => {
    console.log("i want to open this profile");
  };
  function toggleOverlay() {
    showOverlay = !showOverlay;
  }
  function toggleComments() {
    showComments = !showComments;
  }
</script>

<main>
  {#if showOverlay}
    <PostOverlay on:close={toggleOverlay} />
  {/if}

  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div class="createPost" on:click={toggleOverlay}>Create new Post..</div>
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div class="singlePost" >
      <div class="userInfo" on:click={openProfile()}>
        <img src={$userInfo.avatar} alt="user avatar" />
        <p class="username">{$userInfo.firstName} {$userInfo.lastName}</p>
        <p class="createdAt">post.CreatedAt</p>
      </div>

    <div class="postContent" on:click={toggleComments}>
      This is post, i like turtles and please jõuludeks mulle "kolm kotti täis
      viiesajaseid" Click for comments
    </div>
    {#if showComments}
      <div in:slide class="addComment">
        <textarea placeholder="Comment post.."></textarea>
        <div class="postButtons">
          <Button type="secondary">Comment</Button>
          <ImageToComment
            inputIDProp="commentImage"
            fakeInputText="Add Image"
          />
        </div>
      </div>
      <div class="comment">This is comment</div>
      <div class="comment">This is comment</div>
      <div class="comment">This is comment</div>
    {/if}
  </div>
</main>

<style>
  main {
    padding: 4px;
    border-radius: 8px;
    border: solid 1px #333;
  }

  div {
    border: solid 1px #333;
    padding: 8px;
    border-radius: 8px;
  }

  .singlePost {
    display: grid;
    grid-template-columns: 2;
    grid-template-rows: 2;
  }
  
  /* post creator stuff */

  .userInfo {
    display: grid;
    align-items: center;
    grid-column: 1;
  }

  .username {
    grid-row: 1;
  }

  img {
    grid-row: 2;
    border-radius: 50px;
    max-width: 90px;
  }
  .createdAt {
    grid-row: 3;
  }

  /* post creator stuff end  */

  .postContent {
    grid-column: 2;
    display: grid;
  }

  .postButtons {
    padding: 0;
    margin-top: 0.5em;
    border: none;
  }

  .createPost {
    display: flex;
    flex-direction: row;
    color: #555;
    border-color: greenyellow;
    padding: 8px;
    margin: 4px;
  }

  .createPost:hover {
    cursor: pointer;
  }

  .addComment {
    display: grid;
    /* grid-row: 2; */
  }

  /* .addComment textarea {
    margin-right: 4px;
    flex-grow: 1;
    border-radius: 8px;
    min-height: 80px;
    margin-bottom: 0;
    padding: 4px;
    text-align: left;
    vertical-align: top;
    resize: none;
    overflow: auto;
  } */

  .comment {
    margin: 2px;
  }
</style>
