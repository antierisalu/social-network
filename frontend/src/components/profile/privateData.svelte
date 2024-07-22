<script>
  import Matrix from "../../shared/matrix.svelte";
  import { slide, fade } from "svelte/transition";
  import {
    userProfileData,
    isEditingProfile,
    newAboutMeStore,
  } from "../../stores";

  // export let newAboutMe = ''
  import Followers from "./followers.svelte";
  import AllPostsOverlay from "./allPostsOverlay.svelte";

  let x;
  let y;

  $: user = $userProfileData;

  let showOverlay = false;
  let showPostOverlay = false;
  let overlayInfo = [];
  function followOverlay(n, event) {
    x = event.clientX - window.innerWidth / 2;
    y = event.clientY - window.innerHeight / 2;
    console.log(x, y);
    showOverlay = true;
    if (n === 1) {
      overlayInfo = user.followers;
      console.log(n, user.followers);
    } else {
      overlayInfo = user.following;
      console.log(n, user.following);
    }
  }

  export let followerCount;

  function handleAboutMeChange() {}

  function toggleOverlay() {
    showOverlay = !showOverlay;
  }

  function togglePostOverlay() {
    console.log(user)
    showPostOverlay = !showPostOverlay;
  }
</script>

{#if showOverlay && overlayInfo}
  <Followers on:close={toggleOverlay} followers={overlayInfo} {x} {y} />
{/if}

{#if showPostOverlay}
  <AllPostsOverlay on:close={togglePostOverlay} posts={user.posts} />
{/if}

<div class="PrivateData" in:slide>
  <label for="birthday">Birthday</label>
  <div class="birthday">{user.dateOfBirth.String}</div>
  {#if user.aboutMe.String && !$isEditingProfile}
    <label for="aboutMe">About me</label>
    <div in:fade class="aboutMe">{user.aboutMe.String}</div>
  {:else if $isEditingProfile}
    <label in:fade for="aboutMe">About me</label>
    <input
      in:fade
      type="text"
      class="editProfileText"
      bind:value={$newAboutMeStore}
      on:input={handleAboutMeChange}
    />
  {/if}
  <div class="follow">
    <div>
      <label for="followers">Followers</label>
      <div>
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <div class={followerCount > 0 ?"followers":"zero"} on:click={() => followOverlay(1, event)}>
          {followerCount}
        </div>
      </div>
    </div>
    <div>
      <label for="followers">Following</label>
      <div>
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <div class={followerCount > 0 ?"following":"zero"} on:click={() => followOverlay(0, event)}>
          {user.following ? user.following.length : 0}
        </div>
      </div>
    </div>
  </div>
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div class="userPostLabels">
    <label for activity>{user.posts ? "Latest posts" : "User has no posts"}</label>
    {#if user.posts}
    <u on:click={togglePostOverlay}>See all Posts</u>
    {/if}
  </div>
  {#if user.posts === null}
    <Matrix />
  {:else}
    <div class="activity">
      {#each user.posts.slice(0, 5) as post}
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <div class="userPost" on:click={togglePostOverlay}>
          {post.content.slice(0, 30)}
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  label {
    padding: 8px;
    font-weight: bold;
  }
  .userPost {
    border: solid 1px yellowgreen;
    border-radius: 8px;
    padding: 4px;
    margin: 4px 0;
    cursor: pointer;
  }
  .userPostLabels {
    display: flex;
    flex-direction: row;
    justify-content: space-evenly;
  }

  u {
    padding: 8px;
    cursor: grabbing;
  }

  .follow {
    padding: 0;
    display: flex;
    justify-content: center;
    max-height: 200px;
    overflow-y: none;
    scrollbar-width: thin;
    scrollbar-color: greenyellow #011;
  }

  .aboutMe,
  .activity,
  .birthday,
  .zero,
  .following,
  .followers {
    font-size: small;
    max-height: 200px;
    overflow: auto;
    border: solid 1px #333;
    border-radius: 6px;
    text-align: center;
    padding: 8px;
  }

  .zero,
  .following,
  .followers {
    margin: 4px;
    padding: 4px auto;
  }

  .following:hover,
  .followers:hover {
    color: white;
    border-color: white;
    cursor: pointer;
  }

  .activity {
    max-height: 500px;
  }

  .editProfileText {
    width: 100%;
    text-align: center;
    border-color: greenyellow;
    margin: 0;
    padding: 8px;
  }
</style>
