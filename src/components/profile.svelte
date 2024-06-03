<script>
import Button from "../shared/button.svelte";
</script>

<main>
    <div class="userContainer">
        <div class="name">FirstName LastName</div>
        <div class="avatar">avatar</div>
        <div class="buttons">
            <!-- {#if followingUser }  -->
                <Button>unFollow</Button>
                <!-- {:else} -->
                <Button type="secondary" w84={true}>follow</Button>
            <!-- {/if} -->
            <Button type="secondary" w84={true}>Chat</Button>
        </div>
        
        <div class="birthday">Birthday</div>
        <div class="aboutme">About me</div>
        <div class="posts">User activity
            <li>1</li>
            <li>1</li>
            <li>1</li>
            <li>1</li>
            <li>1</li>
        </div>
    </div>
</main>

<style>

    .userContainer {
        display: flex;
        flex-direction: column;
    }

    .posts {
        display: flex;
        flex-direction: column;
        border: solid 1px red;
        bottom: 0;
    }
    

</style>