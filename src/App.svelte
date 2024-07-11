<script>
  import Header from "./components/structure/header.svelte";
  import Login from "./components/auth/login.svelte";
  import Mainpage from "./components/structure/mainpage.svelte";
  import { loggedIn, userInfo, userProfileData } from "./stores";
  import { onMount } from "svelte";
  import { fetchUsers, fetchNotifications } from "./utils";

  let isMounted = false; //et login ei flashiks refreshi ajal

  onMount(async () => {
    try {
      const response = await fetch("/session");
      if (!response.ok) {
        throw new Error("Network response was not ok: " + response.statusText);
      }
      const data = await response.json();
      userInfo.set(data); //set global userInfo for components to access all user info
      userProfileData.set(data); //set global userInfo for components to access all user info
      // console.log($userInfo);
      loggedIn.set(true);
      fetchUsers();
      fetchNotifications();
    } catch (error) {
      console.error("Error fetching session:", error.message);
    } finally {
      isMounted = true;
    }
  });
</script>

<Header />
{#if isMounted}
  <main>
    <svelte:component this={$loggedIn ? Mainpage : Login}></svelte:component>
  </main>
{/if}

<style>
  main {
    text-align: center;
    padding: 1em;
    max-width: 240px;
    margin: 0 auto;
    max-width: 90vh;
  }

  @media (min-width: 640px) {
    main {
      max-width: none;
    }
  }
</style>
