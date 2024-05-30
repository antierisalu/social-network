<script>
  import Footer from "./components/footer.svelte";
  import Header from "./components/header.svelte";
  import Login from "./components/login.svelte";
  import Mainpage from "./components/mainpage.svelte";

  let loggedIn = false;
  function handleLogin(event) {
    console.log(event.detail.loginStatus);
    loggedIn = event.detail.loginStatus;
  }
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
