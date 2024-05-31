<script>
    // import Login from "./login.svelte";
    import Button from "../shared/button.svelte";

</script>

<header>
    {#if}
    <div class="leftside">
        <Button type="secondary" inverse={true}>Chat</Button>
        <Button type="secondary" inverse={true}>Groups</Button>
    </div>
    {/if}
    <h3>Choice is an illusion</h3>
    {#if}
    <div class="rightside">
        <Button type="secondary" flat={true} inverse={true}>Profile/Username</Button>
        <Button type="primary">Logout</Button>
    </div>
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
