<script>
import Button from "../shared/button.svelte";
</script>

<main>
    <div class="userContainer">
        <div class="name">FirstName LastName</div>
        <div class="avatar"></div>
        <div class="buttons">
            <!-- {#if followingUser }  -->
                <Button id="unFollowBtn">unFollow</Button>
                <!-- {:else} -->
                <Button type="secondary" w84={true} id="FollowBtn">follow</Button>
            <!-- {/if} -->
            <Button type="secondary" w84={true} id="chatBtn">Chat</Button>
        </div>
        <div class="birthdayContainer">
            <label for="birthday">Birthday: </label>
            <div class="birthday"> 04-05-2024</div>
        </div>
        <label for="aboutme">About me:</label>
        <div class="aboutme">Every once in a while something amazing comes around, so here i am :)))</div>
        <div class="posts">
            <li>Accepted follower request</li>
            <li>Posted to "Märgatud Viljandis"</li>
            <li>Commented Albert Einstein post</li>
            <li>Created Group "Hack the Planet"</li>
            <li>Invited you to the group "DIY"</li>
            <li>Invited you to the group "Big fat dark marketplace for car keys"</li>
            <li>Created Event "Jaanid 2024"</li>
        </div>
    </div>
</main>

<style>
    div {
        padding-bottom: 8px;
    }

    main {
        padding: 4px;
        color: rgba(172, 255, 47, 0.616);
        display: flex;
        flex-direction: column;
        text-align: center;
    }

    .avatar {
        height: 278px;
    }

    .birthdayContainer {
        display: flex;
        justify-content: space-around;

    }

    .posts {
        display: flex;
        flex-direction: column;
        border: solid 1px red;
        bottom: 0;
    }
    

</style>