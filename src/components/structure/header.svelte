<script>
  import Button from "../../shared/button.svelte";
  import { loggedIn, activeTab, userProfileData, userInfo } from "../../stores";
  import { getGroups } from "../../utils";
  import { blur } from "svelte/transition";

  function logout() {
    loggedIn.set(false);
    document.cookie = `sessionToken=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
  }
</script>

<header>
  {#if $loggedIn}
    <div in:blur class="leftside">
      <Button
        type="secondary"
        w84
        w120={false}
        inverse
        on:click={() => activeTab.set("Notifications")}>🔔</Button
      >
      <Button
        type="secondary"
        w84
        w120={false}
        inverse
        on:click={() => {
          activeTab.set("Groups");
        }}>Groups</Button
      >
      <Button
        type="secondary"
        w84
        w120={false}
        inverse
        on:click={() => {
          activeTab.set("Profile");
          userProfileData.set($userInfo);
        }}>Profile</Button
      >
    </div>
    <h2><a href="/">🏠</a></h2>
    <div in:blur class="rightside">
      <Button type="primary" w84 w120={false} inverse on:click={() => logout()}
        >LogOut</Button
      >
    </div>
  {/if}
</header>

<style>
  header {
    z-index: 10;
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
    right: 8px;
    top: 12px;
  }

  .leftside {
    position: absolute;
    left: 20px;
    top: 12px;
  }
</style>
