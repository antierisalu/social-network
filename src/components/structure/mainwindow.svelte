<script>
  import Button from "../../shared/button.svelte";
  import ImageToPost from "../../shared/imagePreview.svelte";
  import ImageToComment from "../../shared/imagePreview.svelte";

  let showComments = false;

    function toggleComments() {
    showComments = !showComments;
    }


</script>

<main>
    <div class="createPost">
        <textarea placeholder="New post here.."></textarea>
        <div class="postButtons">
            <Button type="secondary">Post</Button>
            <ImageToPost fakeInputText='Add Image'/>  
        </div>
    </div>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="post" >
        <div class="postContent" on:click={toggleComments}>This is post, i like turtles and please jõuludeks mulle "kolm kotti täis viiesajaseid"</div>
        {#if showComments}
          <div class="addComment">
            <textarea placeholder="Comment post.."></textarea>
            <div class="postButtons">
              <Button type="secondary">Comment</Button>
              <ImageToComment fakeInputText='Add Image' />
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
        border: solid 1px #333
    }

    div {
        border: solid 1px #333;
        padding: 8px;
        border-radius: 8px;
    }
    .postButtons {
        display: flex;
        flex-direction: column;
        padding: 0;
        margin-top:0.5em;
        border:none
    }

    .createPost, .addComment {
        display: flex;
        align-items: center;
    }

    .createPost textarea, .addComment textarea {
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
    }

    .comment {
        margin: 2px;
    }
    .post {
        text-align: left;
        margin: 4px
    }
</style>