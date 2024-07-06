<script>
  import Button from "../../shared/button.svelte";
  import { slide } from "svelte/transition";
  import Comment from "./comments.svelte";
  import PostOverlay from "./createPost.svelte";
  import ImageToComment from "../../shared/imagePreview.svelte";
  import { uploadImageStore } from "../../stores";
  import {
    getUserDetails,
    getPosts,
    selectUser,
    getComments,
  } from "../../utils";
  import { writable } from "svelte/store";

  let showOverlay = false;
  let commentsVisibility = writable([]);
  let comments = [];
  let newCommentContent = "";
  export let posts;
  export let allowCreate = true;
  
  let uploadImage;
  uploadImageStore.subscribe((value) => {
    uploadImage = value;
  });

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
    comments = posts[index].comments;
  }

  async function sendComment(postID, post) {
    const response = await fetch("http://localhost:8080/newComment", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ postID, content: newCommentContent }),
    });
    const comment = await response.json();
    if (!response.ok) {
      console.error("Failed to add comment");
    }
    await uploadImage({ comment: comment.id }).catch((error) => {
      console.error("Error uploading the image:", error);
    });
    newCommentContent = "";

    post.comments = await getComments(postID);
    comments = post.comments;
    console.log(comments, post.comments);
  }

  $: if (posts) commentsVisibility.set(Array(posts.length).fill(false));
</script>

<main>
  {#if showOverlay}
    <PostOverlay on:close={toggleOverlay} />
  {/if}

  <!-- svelte-ignore a11y-click-events-have-key-events -->
  {#if allowCreate}
    <div class="createPost" on:click={toggleOverlay}>Create new Post..</div>
  {/if}

  {#if posts}
    <div class="postsFeed">
      {#each posts as post, index}
        {#await Promise.resolve(getUserDetails(post.userID)) then user}
          {#if user}
            <div class="singlePost">
              <!-- svelte-ignore a11y-click-events-have-key-events -->
              <div class="userInfo" on:click={() => selectUser(user.ID)}>
                <p class="postCreator">{user.FirstName} {user.LastName}</p>
                <p class="postCreatorAvatar">
                  <img
                    class="postCreatorAvatar"
                    src={user.Avatar}
                    alt="user avatar"
                  />
                </p>
                <p class="postCreatedAt">{post.createdAt}</p>
              </div>
              <!-- svelte-ignore a11y-click-events-have-key-events -->
              <div class="postContent" on:click={() => toggleComments(index)}>
                <div class="content">
                  <div>{post.content}</div>
                  <p><img src={post.img} alt={post.img} /></p>
                </div>
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
                      on:click={() => {
                        if (newCommentContent !== "")
                          sendComment(post.id, post);
                        else alert("Comment cannot be empty");
                      }}>Comment</Button
                    >
                    <ImageToComment
                      inputIDProp="commentImage"
                      fakeInputText="Add Image"
                    />
                  </div>
                </div>
                {#if post.comments}
                  <div class="comments">
                    {#each comments as comment}
                      <Comment {comment} />
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
    padding: 4px;
    border-radius: 8px;
  }

  p {
    margin: 0;
  }

  img {
    max-height: 300px;
    max-width: 300px;
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

  .postCreatorAvatar {
    margin: 4px 0;
    max-height: 100px;
    border-radius: 10px;
  }

  textarea {
    width: 100%;
    min-height: 60px;
    resize: vertical;
  }

  .singlePost {
    border: solid 1px greenyellow;
    display: grid;
    grid-auto-columns: 1fr;
    grid-template-columns: 0.3fr 1.5fr;
    grid-template-rows: 0.5fr 0.5fr minmax(0 3fr);
    gap: 0px 0px;
    margin: 4px 0;
    grid-template-areas:
      "userInfo postContent"
      "addComment addComment"
      "comments comments";
  }
  .postContent {
    grid-area: postContent;
    /* margin: 0 8px; */
    border: solid 1px #333;
  }
  .userInfo {
    grid-area: userInfo;
    cursor: pointer;
  }
  .addComment {
    display: grid;
    grid-area: addComment;
    grid-template-columns: auto 150px;
  }

  .addComment > textarea {
    grid-column: 1;
  }

  .commentButtons {
    display: flex;
    /* height: 80px; This cannot be used, when adding image, it goes into another post but it needs to be able to grow */
    flex-direction: column;
    grid-column: 2;
  }

  .comments {
    grid-area: comments;
  }
</style>
