<script>
  import Chat from "../chat/chat.svelte";
  import { fade, slide } from "svelte/transition";
  import Footer from "./footer.svelte";
  import Profile from "../profile/profile.svelte";
  import Notifications from "../notifications/notifications.svelte";
  import Groups from "../groups/groups.svelte";
  import MainWindow from "./mainwindow.svelte";
  import { activeTab, userInfo } from "../../stores";
  import { connect, sendMessage, messages } from "../../websocket";
  import { onMount } from "svelte";
  import Button from "../../shared/button.svelte";
  import UserSearch from "../profile/searchBar.svelte"

  onMount(() => {
    console.log("connecting ws", $userInfo);
    connect($userInfo.firstName);
  });

</script>

<main in:fade>
  
  <div id="leftSidebar" in:fade>
    {#if $activeTab === "Profile"}
    <UserSearch />
      <div in:fade><Profile /></div>
    {:else if $activeTab === "Groups"}
      <div in:fade><Groups /></div>
    {:else}
      <div in:fade><Notifications /></div>
    {/if}
  </div>

  <div id="mainWindow">
    <Button inverse={true} on:click={() => sendMessage(JSON.stringify({ type: "ping", data: "ping" }))}>send</Button>
    <!--if groups
      else posts
      else blablabla-->
    <MainWindow />
  </div>
  <div id="rightSidebar">
    <Chat />
  </div>
  <div id="footer">
    <Footer />
  </div>
</main>

<style>
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



  #footer {
    grid-column: 2/4;
  }
  #leftSidebar {
    grid-row: 1/3;
  }
</style>
