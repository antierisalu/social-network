<script>
    import { onMount } from 'svelte';
    import Button from "../shared/button.svelte";

    let notOwnPage

    let user = {
        id: '',
        email: '',
        firstName: '',
        lastName: '',
        dateOfBirth: '',
        avatar: '',
        nickName: '',
        aboutMe: '',
    };

    onMount(async () => {
        try {
            const response = await fetch('http://localhost:8080/session'); // Replace with your actual endpoint
            const data = await response.json();
            user = data;
            user.avatar = "./avatars/av.png"
        } catch (error) {
            console.error('Error fetching user data:', error);
        }
    });</script>

<main>
    
    <div class="userContainer">
        <div class="name">{user.firstName} {user.lastName} {#if user.nickName}({user.nickName}){/if}</div>
        <div class="avatar">
            <img src={user.avatar} border="0" alt="avatar" />
        </div>
        {#if notOwnPage}
        <div class="buttons">
            <!-- {#if followingUser }  -->
                <Button id="unFollowBtn">unFollow</Button>
                <!-- {:else} -->
                <Button type="secondary" w84={true} id="followBtn">Follow</Button>
            <!-- {/if} -->
            <Button type="secondary" inverse={true} w84={true} id="chatBtn">Chat</Button>
        </div>
        {/if}
            <div class="birthday">Birthday: {user.dateOfBirth}</div>
            {#if user.aboutMe}
        <label for="aboutMe">About me:</label>
        <div class="aboutMe">{user.aboutMe}</div>
            {/if}
        <label for activity>Latest posts</label>
        <div class="activity">
            <div>Posted to "Märgatud Viljandis"</div>
            <div>Posted to "Märgatud Viljandis"</div>
            <div>Posted to "Märgatud Viljandis"</div>
        </div>
    </div>
</main>

<style>


main {
        padding: 8px;
        color: rgba(172, 255, 47, 0.616);
        display: flex;
        flex-direction: column;
        font-size: small;
    }

    div {
        padding-bottom: 8px;
    }

    img {
        max-width: 280px;
    }
    label {
        padding: 8px;
    }

    .aboutMe, .activity{
        font-size: small;
        max-height: 200px;
        overflow: auto;
        border: solid 1px greenyellow;
        border-radius: 6px;
        text-align: left;
        padding: 8px;
    }

    .activity {
        max-height: 500px;
    }
    

</style>