<script>
  import Chat from "./chat.svelte";
  import { fade, slide } from "svelte/transition";
  import Footer from "./footer.svelte";
  import Profile from "./profile.svelte";
  import Notifications from "./notifications.svelte";
  import Groups from "./groups.svelte";
  import MainWindow from "./mainwindow.svelte";
  import { activeTab } from "../stores";
  import { connect, sendMessage, messages } from "../websocket";
  import { onMount } from "svelte";
  import Button from "../shared/button.svelte";

  onMount(() => {
    console.log("connecting ws");
    connect();
  });

  $: console.log($activeTab);
</script>

<main in:fade>
  <Button inverse={true}
    on:click={() => sendMessage(JSON.stringify({ type: "ping", data: "ping" }))}
    >send</Button
  >
  <div id="leftSidebar" in:fade>

    {#if $activeTab === "Profile"}
    <div class="searchUsers">
      <input type="text" placeholder="Search users">
    </div>
      <div in:fade><Profile /></div>
    {:else if $activeTab === "Groups"}
      <div in:fade><Groups /></div>
    {:else}
      <div in:fade><Notifications /></div>
    {/if}
  </div>
  <div id="mainWindow">
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

  input{
        margin: 0;
        border: none;
        width: 100%;
        height: 100%;
    }
  
    .searchUsers {
      margin-bottom: 8px;
    }

  #footer {
    grid-column: 2/4;
  }
  #leftSidebar {
    grid-row: 1/3;
  }
</style>
