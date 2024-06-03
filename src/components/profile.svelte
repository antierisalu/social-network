<script>
import Button from "../shared/button.svelte";
</script>

<main>
    <div class="userContainer">
        <div class="name">FirstName LastName</div>
        <div class="avatar">
            <img src="../avatars/av.png" border="0" alt="av.png" />
        </div>
        <!-- {#if notOwnPage} -->
        <div class="buttons">
            <!-- {#if followingUser }  -->
                <Button id="unFollowBtn">unFollow</Button>
                <!-- {:else} -->
                <Button type="secondary" w84={true} id="followBtn">Follow</Button>
            <!-- {/if} -->
            <Button type="secondary" inverse={true} w84={true} id="chatBtn">Chat</Button>
        </div>
        <!-- {/if} -->
            <div class="birthday">Birthday: 04-05-2000</div>
        <label for="aboutMe">About me:</label>
        <div class="aboutMe">Every once in a while something amazing comes around, so here i am :)))Every once in a while something amazing comes around, so here i am :)))</div>
        <div class="activity">
            <div>Accepted follower request</div>
            <div>Posted to "Märgatud Viljandis"</div>
            <div>Commented Albert Einstein post</div>
            <div>Created Group "Hack the Planet"</div>
            <div>Invited you to the group "DIY"</div>
            <div>Invited you to the group "Big fat dark marketplace for car keys"</div>
            <div>Created Event "Jaanid 2024"</div>
            <div>Accepted follower request</div>
            <div>Posted to "Märgatud Viljandis"</div>
            <div>Commented Albert Einstein post</div>
            <div>Created Group "Hack the Planet"</div>
            <div>Invited you to the group "DIY"</div>
            <div>Invited you to the group "Big fat dark marketplace for car keys"</div>
            <div>Created Event "Jaanid 2024"</div>
            <div>Accepted follower request</div>
            <div>Posted to "Märgatud Viljandis"</div>
            <div>Commented Albert Einstein post</div>
            <div>Created Group "Hack the Planet"</div>
            <div>Invited you to the group "DIY"</div>
            <div>Invited you to the group "Big fat dark marketplace for car keys"</div>
            <div>Created Event "Jaanid 2024"</div>
            <div>Accepted follower request</div>
            <div>Posted to "Märgatud Viljandis"</div>
            <div>Commented Albert Einstein post</div>
            <div>Created Group "Hack the Planet"</div>
            <div>Invited you to the group "DIY"</div>
            <div>Commented Albert Einstein post</div>
            <div>Created Group "Hack the Planet"</div>
            <div>Invited you to the group "DIY"</div>
            <div>Commented Albert Einstein post</div>
            <div>Created Group "Hack the Planet"</div>
            <div>Invited you to the group "DIY"</div>
            <div>Commented Albert Einstein post</div>
            <div>Created Group "Hack the Planet"</div>
            <div>Invited you to the group "DIY"</div>
            <div>Invited you to the group "Big fat dark marketplace for car keys"</div>
            <div>Created Event "Jaanid 2024"</div>
        </div>
    </div>
</main>

<style>

    main {
        padding: 4px;
        color: rgba(172, 255, 47, 0.616);
        display: flex;
        flex-direction: column;
        text-align: center;
        max-height: 90vh;
    }

    div {
        padding-bottom: 8px;
    }

    img {
        max-width: 280px;
    }

    .aboutMe {
        max-height: 200px;
        overflow: auto;
    }

    .activity {
        overflow: auto;
        /* max-height: 600px;  */
        height: 55vh; /*i dont know about that*/
        padding: 4px;
        display: flex;
        flex-direction: column;
        border-radius: 6px;
        border-color: greenyellow;
        border: solid 1px greenyellow;
        text-align: left;
    }
    

</style>