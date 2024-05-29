<script>
    import Button from "../shared/button.svelte";
    import { createEventDispatcher } from "svelte";
    const dispatch = createEventDispatcher();


    export let loggedIn = false;

    export function setLoggedIn() {
        dispatch('login',{
            loginStatus: loggedIn
        })
    }

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

</script>

<div>
    <form on:submit|preventDefault>
        <input type="text" placeholder="username" bind:value={user} required>
        <input type="password" placeholder="password" required bind:value={password}>
        <Button type="secondary" on:click={fetchData}>Login</Button>
        
    </form>
</div>

<style>

    input{
        border-radius: 6px;
        padding: 8px 12px;
    }
</style>