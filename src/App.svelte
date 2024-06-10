<script>
  import Header from "./components/structure/header.svelte";
  import Login from "./components/auth/login.svelte";
  import Mainpage from "./components/structure/mainpage.svelte";
  import { loggedIn, userInfo, allUsers } from "./stores";
  import { onMount } from "svelte";

  let isMounted = false; //et login ei flashiks refreshi ajal

  const fetchUsers = async () => {
        const response = await fetch('http://localhost:8080/allusers');
        if(response.ok) {
            const fetchedUsers = await response.json();
            $allUsers = [...fetchedUsers];
        } else {
            console.error('Error fetching users:', response.status);
        }
    };

  onMount(async () => {
    try {
      const response = await fetch("/session");
      if (!response.ok) {
        throw new Error("Network response was not ok: " + response.statusText);
      }
      const data = await response.json();
      userInfo.set(data);//set global userInfo for components to access all user info
      loggedIn.set(true);
      fetchUsers()
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
  }

  @media (min-width: 640px) {
    main {
      max-width: none;
    }
  }
</style>
