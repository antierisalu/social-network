<script>
  import Button from "../../shared/button.svelte";
  import { slide } from "svelte/transition";
  import PostOverlay from "./createPost.svelte";
  import ImageToComment from "../../shared/imagePreview.svelte";
  import { posts, userInfo } from "../../stores";
  import { getUserDetails, getPosts, selectUser } from "../../utils";
  import { writable } from "svelte/store";

  let showOverlay = false;
  let commentsVisibility = writable([]);
  let newCommentContent = "";

  const openProfile = (userID) => {
    console.log(`i want to open this profile ${userID}`);
  };

  export function toggleOverlay() {
    showOverlay = !showOverlay;
    if (!showOverlay) {
      getPosts();
    }
  }

  function toggleComments(index) {
    commentsVisibility.update((current) => {
      const updated = current.map((visible, i) =>
        i === index ? !visible : false
      );
      return updated;
    });
  }

  async function sendComment(postID) {
    const response = await fetch("http://localhost:8080/newComment", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ postID, content: newCommentContent }),
    });

    if (!response.ok) {
      console.error("Failed to add comment");
    }
    newCommentContent = "";
  }

  $: if ($posts) commentsVisibility.set(Array($posts.length).fill(false));
</script>

<main>
  {#if showOverlay}
    <PostOverlay on:close={toggleOverlay} />
  {/if}

  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div class="createPost" on:click={toggleOverlay}>Create new Post..</div>

  {#if $posts}
    <div class="postsFeed">
      {#each $posts as post, index}
        {#await Promise.resolve(getUserDetails(post.userID)) then user}
          {#if user}
            <div class="singlePost">
              <!-- svelte-ignore a11y-click-events-have-key-events -->
              <div class="userInfo" on:click={() => selectUser(user.ID)}>
                <p class="postCreator">{user.FirstName} {user.LastName}</p>
                <p class="postCreatorAvatar">
                  <img src={user.Avatar} alt="user avatar" />
                </p>
                <p class="postCreatedAt">{post.createdAt}</p>
              </div>
              <!-- svelte-ignore a11y-click-events-have-key-events -->
              <div class="postContent" on:click={() => toggleComments(index)}>
                {@html post.content}
              </div>

              {#if $commentsVisibility[index]}
                <div in:slide class="addComment">
                  <textarea
                    bind:value={newCommentContent}
                    placeholder="Comment post.."
                  ></textarea>
                  <div class="commentButtons">
                    <Button
                      type="secondary"
                      on:click={() => sendComment(post.id)}>Comment</Button
                    >
                    <ImageToComment
                      inputIDProp="commentImage"
                      fakeInputText="Add Image"
                    />
                  </div>
                </div>
                {#if post.comments}
                  <div class="comments">
                    {#each post.comments as comment}
                      <div class="singleComment">
                        <!-- svelte-ignore a11y-click-events-have-key-events -->
                        <div
                          class="userInfo"
                          on:click={() => openProfile(comment.userID)}
                        >
                          <p class="commentCreator">
                            {comment.user.firstName}
                            {comment.user.lastName}
                          </p>
                          <p class="commentCreatorAvatar">
                            <img src={comment.user.avatar} alt="user avatar" />
                          </p>
                          <p class="commentCreatedAt">{comment.createdAt}</p>
                        </div>
                        <div class="commentContent">{comment.content}</div>
                      </div>
                    {/each}
                  </div>
                {/if}
              {/if}
            </div>
          {/if}
        {/await}
      {/each}
    </div>
  {/if}
</main>

<style>
  main {
    display: flex;
    flex-direction: column;
    padding: 4px;
    border-radius: 8px;
    border: solid 1px #333;
    height: 85vh;
    overflow-y: scroll;
    scrollbar-width: thin;
    scrollbar-color: greenyellow #011;
  }

  div {
    padding: 8px;
    border-radius: 8px;
  }

  p {
    margin: 0;
  }

  .createPost {
    display: flex;
    flex-direction: row;
    color: #555;
    border: solid 1px greenyellow;
    padding: 8px;
    margin: 4px;
  }
  .createPost:hover {
    cursor: pointer;
  }

  img {
    padding: 4px 0;
    max-height: 60px;
  }

  textarea {
    width: 100%;
    min-height: 100px;
    resize: vertical;
  }

  .singlePost {
    border: solid 1px #333;
    display: grid;
    grid-auto-columns: 1fr;
    grid-template-columns: 0.3fr 1.5fr;
    grid-template-rows: 0.5fr 0.5fr minmax(0 3fr);
    gap: 0px 0px;
    grid-template-areas:
      "userInfo postContent"
      "addComment addComment"
      "comments comments";
  }
  .postContent {
    grid-area: postContent;
    margin: 0 8px;
    border: solid 1px #333;
  }
  .userInfo {
    grid-area: userInfo;
    cursor: pointer;
  }
  .addComment {
    grid-area: addComment;
  }
  .comments {
    grid-area: comments;
  }

  .singleComment {
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

  .commentContent {
    grid-area: commentContent;
    margin: 0 8px;
    border: solid 1px #333;
  }
</style>
