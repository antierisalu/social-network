<script>
  import Footer from "./components/footer.svelte";
  import Header from "./components/header.svelte";
  import Login from "./components/login.svelte";
  import Mainpage from "./components/mainpage.svelte";
  import { onMount } from "svelte";

  let loggedIn = false;
  function handleLogin(event) {
    loggedIn = event.detail.loginStatus;

    if (!loggedIn) {
      document.cookie =
        "sessionToken=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    }
  }

  onMount(async () => {
    try {
      const response = await fetch("/session");
      if (!response.ok) {
        throw new Error("Network response was not ok: " + response.statusText);
      }
      const data = await response.json();
      loggedIn = true;
    } catch (error) {
      console.error("Error fetching session:", error.message);
    }
  });
</script>

<Header />
<main>
  <svelte:component
    this={loggedIn ? Mainpage : Login}
    on:login={handleLogin}
    on:logout={handleLogin}
  ></svelte:component>
</main>

<Footer />

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
