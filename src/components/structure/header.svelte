<script>
  import Button from "../../shared/button.svelte";
  import { loggedIn, activeTab, userProfileData, userInfo } from "../../stores";
  import { blur } from "svelte/transition";
  import { sendMessage } from "../../websocket";


  function logout() {
    sendMessage(JSON.stringify({ type: "logout", data: "", username:$userInfo.email }));
    loggedIn.set(false);
    document.cookie = `sessionToken=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
  }
</script>

<header>
  {#if $loggedIn}
    <div in:blur class="leftside">
      <Button
        type="secondary"
        w84={true}
        inverse={true}
        on:click={() => activeTab.set("Notifications")}>🔔</Button
      >
      <Button
        type="secondary"
        w84={true}
        inverse={true}
        on:click={() => activeTab.set("Groups")}>Groups</Button
      >
      <Button
        type="secondary"
        w84={true}
        inverse={true}
        on:click={() => {
          activeTab.set("Profile");
          userProfileData.set($userInfo);
        }}>Profile</Button
      >
    </div>
    <h2>Choi is an illusion</h2>
    <div in:blur class="rightside">
      <Button type="primary" w84={true} inverse={true} on:click={() => logout()}
        >LogOut</Button
      >
    </div>
  {:else}
    <h2>Choie is an illusion</h2>
  {/if}
</header>

<style>
  header {
    position: sticky;
    top: 0;
    background: #011;
    padding: 12px;
  }

  h2 {
    margin: 0 auto;
    text-align: center;
    color: #e5e5e5;
  }

  .rightside {
    position: absolute;
    right: 12px;
    top: 12px;
  }

  .leftside {
    position: absolute;
    left: 12px;
    top: 12px;
  }
</style>
