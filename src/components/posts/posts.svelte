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
  // function toggleComments() {
  //   showComments = !showComments;
  // }



  const posts = [
    {
        content: "Hello, world! This is my first post.",
        createdBy: "Alice",
        createdAt: "2024-06-19 16:30:00",
        comments: [
            { createdBy: "Eve", createdAt: "2024-06-19 16:35:00", content: "Nice start!" },
            { createdBy: "Bob", createdAt: "2024-06-19 16:40:00", content: "Keep it up!" },
            { createdBy: "Carol", createdAt: "2024-06-19 16:45:00", content: "Looking forward to more." },
            { createdBy: "Nora", createdAt: "2024-06-19 19:35:00", content: "Kyoto is on my bucket list!" },
            { createdBy: "Henry", createdAt: "2024-06-19 19:40:00", content: "Share your favorite spots!" },
            { createdBy: "Isabella", createdAt: "2024-06-19 19:45:00", content: "Beautiful photos!" }
        ]
    },
    {
        content: "Exciting news! 🎉",
        createdBy: "Bob",
        createdAt: "2024-06-19 17:15:00",
        comments: [
            { createdBy: "Alice", createdAt: "2024-06-19 17:20:00", content: "What's the news?" },
            { createdBy: "Eve", createdAt: "2024-06-19 17:25:00", content: "Share the details!" },
            { createdBy: "Dan", createdAt: "2024-06-19 17:30:00", content: "Congratulations!" }
        ]
    },
    {
        content: "Exploring the Hidden Gems of Kyoto 🌸",
        createdBy: "Emily",
        createdAt: "2024-06-19 19:30:00",
        comments: [
            { createdBy: "Nora", createdAt: "2024-06-19 19:35:00", content: "Kyoto is on my bucket list!" },
            { createdBy: "Henry", createdAt: "2024-06-19 19:40:00", content: "Share your favorite spots!" },
            { createdBy: "Isabella", createdAt: "2024-06-19 19:45:00", content: "Beautiful photos!" }
        ]
    },
    {
        content: "Delicious Chocolate Chip Cookies Recipe 🍪",
        createdBy: "Grace",
        createdAt: "2024-06-19 18:00:00",
        comments: [
            { createdBy: "Liam", createdAt: "2024-06-19 18:05:00", content: "I love chocolate chip cookies!" },
            { createdBy: "Sophie", createdAt: "2024-06-19 18:10:00", content: "Can't wait to try this recipe!" },
            { createdBy: "Oliver", createdAt: "2024-06-19 18:15:00", content: "Thanks for sharing!" }
        ]
    },
  ];

  let commentsVisibility = Array(posts.length).fill(false);

  function toggleComments(index) {
    commentsVisibility = commentsVisibility.map((visible, i) => i === index ? !visible : visible);
  }

</script>

<main>
  {#if showOverlay}
    <PostOverlay on:close={toggleOverlay} />
  {/if}

  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div class="createPost" on:click={toggleOverlay}>Create new Post..</div>
  
  {#each posts as post, index }
    <div class="singlePost">
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <div class="userInfo" on:click={openProfile()}>
        <p class="postCreator">{post.createdBy}</p>
        <p class="postCreatorAvatar"><img src="https://i.pravatar.cc/100" alt="user avatar" /></p>
        <p class="postCreatedAt">{post.createdAt}</p>
      </div>
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <div class="postContent" on:click={() => toggleComments(index)}>{post.content}</div>

      {#if commentsVisibility[index]}
        <div in:slide class="addComment">
          <textarea placeholder="Comment post.."></textarea>
          <div class="commentButtons">
            <Button type="secondary">Comment</Button>
            <ImageToComment inputIDProp="commentImage" fakeInputText="Add Image" />
          </div>
        </div>
        <div class="comments">
          {#each post.comments as comment}
            <div class="singleComment">
              <!-- svelte-ignore a11y-click-events-have-key-events -->
              <div class="userInfo" on:click={openProfile()}>
                <p class="commentCreator">{comment.createdBy}</p>
                <p class="commentCreatorAvatar"><img src="https://i.pravatar.cc/100" alt="user avatar" /></p>
                <p class="commentCreatedAt">{comment.createdAt}</p>
              </div>
              <div class="commentContent">{comment.content}</div>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  {/each}
</main>

<style>
  main {
    display: flex;
    flex-direction: column;
    padding: 4px;
    border-radius: 8px;
    border: solid 1px #333;
  }

 div {
    border: solid 1px #333;
    padding: 8px;
    border-radius: 8px;
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

  img {
    max-height: 90px; width: 90px
  }

  textarea {
    width: 100%;
    min-height: 100px;
    resize: vertical;
  }

  .singlePost {
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
.postContent { grid-area: postContent; }
.userInfo { grid-area: userInfo; }
.addComment { grid-area: addComment; }
.comments { grid-area: comments; }

.singleComment {
  display: grid;
  grid-auto-columns: 1fr;
  grid-template-columns: 0.3fr 0.3fr 1.5fr; 
  grid-template-rows: 0.5fr 0.5fr 3fr; 
  grid-template-areas: 
    ". userInfo commentContent"
    ". userInfo commentContent"
    ". userInfo commentContent"; 
}

.commentContent { grid-area: commentContent }



</style>
