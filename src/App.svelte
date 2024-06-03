<script>
  import Header from "./components/header.svelte";
  import Login from "./components/login.svelte";
  import Mainpage from "./components/mainpage.svelte";
  import { loggedIn } from "./stores";
  import { onMount } from "svelte";

  onMount(async () => {
    console.log(loggedIn)
    try {
      const response = await fetch("/session");
      if (!response.ok) {
        throw new Error("Network response was not ok: " + response.statusText);
      }
      const data = await response.json();
      console.log(data);
      loggedIn.set(true);
    } catch (error) {
      console.error("Error fetching session:", error.message);
    }
  });

  loggedIn.subscribe(value => {
    console.log('loggedIn:', value);
  });

</script>

<Header />
<main>
  <svelte:component
    this={$loggedIn ? Mainpage : Login}
  ></svelte:component>
</main>

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
