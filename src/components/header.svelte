<script>
    // import Login from "./login.svelte";
    import Button from "../shared/button.svelte";
    import { loggedIn } from "../stores";

    function logout(){
        loggedIn.set(false)
        document.cookie = document.cookie = `sessionToken=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
    }

</script>

<header>
    {#if $loggedIn}
    <div class="leftside">
        <Button type="secondary" inverse={true}>Chat</Button>
        <Button type="secondary" inverse={true}>Groups</Button>
    </div>
    <h3>Choice is an illusion</h3>
    <div class="rightside">
    <Button type="primary" on:click={()=>logout()}>Log Out</Button>
    <Button type="secondary" flat={true} inverse={true}>Profile/Username</Button>
    </div>

    {:else}
    <h3 style="margin: auto;">Choice is an illusion</h3>
    {/if}
</header>

<style>


    header{
         position: sticky;
         display: flex;
         justify-content: space-between;
         text-align: center;
         background: #011;
         padding: 20px;
    }

    h3 {
        margin: 0;
        text-align: center;
        color: #e5e5e5;
    }

    .rightside {
        right: 0
    }

    .leftside {
        left: 0
    }

</style>
