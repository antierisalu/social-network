<script>

    import Register from "./register.svelte";
    import Button from "../shared/button.svelte";
    import { createEventDispatcher } from "svelte";
    const dispatch = createEventDispatcher();


    export let loggedIn = false;

    export function setLoggedIn() {
        dispatch('login',{
            loginStatus: loggedIn
        })
    }

    $: login = true
    let user;
    let password;

    async function fetchData() {
        let cred = {
            username: user,
            password: password
        }
        console.log(cred)
        try {
            const response = await fetch('/login', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
          body: JSON.stringify(cred)
       });
      
            if (!response.ok) {
            throw new Error('Network response was not ok');
        } else {
            loggedIn = true
            setLoggedIn();
        }

        }catch (error) {
      // Handle any errors that occurred during the fetch
      console.error('Error fetching data:', error.message);
      // You can choose to return null or throw the error further
      return null;
        }
    }

    function switchToRegister() {
    login = !login; // This will show the Register component
  }


</script>

<div class="login">
    {#if login}
    <form on:submit|preventDefault>
        <input type="text" placeholder="E-mail" bind:value={user} required>
        <input type="password" placeholder="Password" required bind:value={password}>
        <Button type="secondary" on:click={fetchData}>Login</Button>
    </form>
    <div class="regBtn">
    <Button type="" btn200px={true} on:click={switchToRegister}>Register Instead</Button>
    </div>
    {:else}

    <Register on:click={switchToRegister} />
    {/if}
</div>

<style>

    .login {
        display: flex;
        flex-direction: column;
        align-items: center;
    }

    form {
        display: flex;
        flex-direction: column;
    }

    input{
        width: 300px;
        border-radius: 6px;
        padding: 8px 12px;
    }

</style>