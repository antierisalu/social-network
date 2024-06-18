<script>
  import Button from "../../shared/button.svelte";
  import ImageToPost from "../../shared/imagePreview.svelte";
  import { slide } from "svelte/transition";
  import PostOverlay from "./createPost.svelte";
  import ImageToComment from "../../shared/imagePreview.svelte";

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
  <div class="singlePost" on:click|once={toggleComments}>
    <div class="postCreator" on:click={openProfile()}>
      <p>user.FirstName user.LastName</p>
      <p>user.Avatar</p>
      <p>post.createdAt</p>
    </div>

    <!-- // ??? -->
    <div class="postContent">
      This is post, i like turtles and please jõuludeks mulle "kolm kotti täis
      viiesajaseid"...
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
    text-align: left;
    margin: 4px;
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

  .createPost,
  .addComment {
    display: flex;
    align-items: center;
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
