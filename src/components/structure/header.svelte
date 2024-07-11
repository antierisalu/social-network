<script>
  import Button from "../../shared/button.svelte";
  import { loggedIn, activeTab, userProfileData, userInfo } from "../../stores";
  import { blur } from "svelte/transition";
  import { notifications } from '../../websocket.js';
  import { onDestroy } from 'svelte';

  let notificationCount = 0;


  const unsubscribe = notifications.subscribe(items => {
        let unseenNotifications = items.filter(notification => !notification.seen);
        notificationCount = unseenNotifications.length;
  });
  onDestroy(unsubscribe);




  function logout() {
    loggedIn.set(false);
    document.cookie = `sessionToken=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
  }
</script>

<header>
  {#if $loggedIn}
    <div in:blur class="leftside">
      <Button
        id="notifbell"
        type="secondary"
        w84={true}
        inverse={true}
        on:click={() => activeTab.set("Notifications")}>🔔 {#if notificationCount > 0}<span class="notif-count">{notificationCount}</span>{/if}</Button
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
    <h2>Choice is an illusion</h2>
    <div in:blur class="rightside">
      <Button type="primary" w84={true} inverse={true} on:click={() => logout()}
        >LogOut</Button
      >
    </div>
  {:else}
    <h2>Choice is an illusion</h2>
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
    right: 12px;
    top: 12px;
  }

  .leftside {
    position: absolute;
    left: 12px;
    top: 12px;
  }
</style>
