<script>
  // import Chat from "../chat/chat.svelte";
  import UserList from "../chat/userList.svelte";
  import { fade } from "svelte/transition";
  // import Footer from "./footer.svelte";
  import Profile from "../profile/profile.svelte";
  import Notifications from "../notifications/notifications.svelte";
  import Groups from "../groups/groups.svelte";
  import Posts from "../posts/posts.svelte";
  import { activeTab, chatTabs, userInfo, allPosts } from "../../stores";
  import { connect } from "../../websocket";
  import { onMount } from "svelte";
  import UserSearch from "../profile/searchBar.svelte";
  import { getPosts } from "../../utils";
  import ChatTabs from "../chat/chatTabs.svelte";

  onMount(() => {
    console.log("connecting ws", $userInfo);
    connect($userInfo.email);
    getPosts();
  });
</script>

<main in:fade>
  <div id="leftSidebar" in:fade>
    {#if $activeTab === "Profile"}
      <div in:fade><UserSearch /></div>
      <div in:fade><Profile /></div>
    {:else if $activeTab === "Groups"}
      <div in:fade><Groups /></div>
    {:else}
      <div in:fade><Notifications /></div>
    {/if}
  </div>

  <div id="mainWindow">
    <Posts posts={$allPosts} />
  </div>
  <div id="rightSidebar" in:fade>
    <UserList />
  </div>
  <div id="bottomChatContainer" in:fade>
    <ChatTabs />
    <!-- <Chatbox /> instances of different user chats will be inside this-->
    <!-- <Chatbox /> -->
  </div>
</main>

<style>
  #bottomChatContainer {
    padding: 0;
    grid-column: 2/3;
    height: 100%;
    max-height: 48px;
    display: flex;
    flex-direction: row;
    justify-content: right;
    align-items: center;
  }
  main {
    display: grid;
    grid-template-columns: 300px auto 220px;
    grid-template-rows: auto 50px;
    height: 90vh;
  }
  div {
    padding: 8px;
    border-radius: 8px;
    border: solid 1px #333;
  }
  #leftSidebar {
    grid-row: 1/3;
  }

  #mainWindow {
    grid-row: 1/2;
  }

  #rightSidebar {
    min-height: 85vh;
    overflow: hidden;
    grid-row: 1/3;
  }
</style>
