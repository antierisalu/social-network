<script>

    import { onMount } from 'svelte';
    import Button from "../shared/button.svelte";
    import Matrix from '../shared/matrix.svelte';


    // For buttons to work
    // let followingUser 
    // let followRequested


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
        posts: '',
    };

    onMount(async () => {
        try {
            const response = await fetch('http://localhost:8080/session'); // Replace with your actual endpoint
            const data = await response.json();
            user = data;
            user.posts = ['123', 'Hello'];
            // user.avatar = ''
        } catch (error) {
            console.error('Error fetching user data:', error);
        }
    });</script>

<main>
    
    <div class="userContainer">
        <div class="name">{user.firstName} {user.lastName}</div>
        {#if user.nickName}
        <label for="name">({user.nickName})</label>
        {/if}
        {#if user.avatar}
            <div class="avatar">
                <img src={user.avatar} border="0" alt="avatar" />
            </div>
        {:else}
            <Matrix />
        {/if}
        {#if notOwnPage}
        <div class="buttons">
            <!-- {#if followingUser }  -->
                <Button id="unFollowBtn">unFollow</Button>
                <!-- {:else if !followingUser && followRequested} -->
                <Button id="unFollowBtn">Cancel request</Button>
                <!-- {:else } -->
                <Button type="secondary" w84={true} id="followBtn">Follow</Button>
            <!-- {/if} -->
            <Button type="secondary" inverse={true} w84={true} id="chatBtn">Chat</Button>
        </div>
        {/if}
            <label for="birthday">Birthday</label>
            <div class="birthday">{user.dateOfBirth}</div>
        {#if user.aboutMe}
            <label for="aboutMe">About me</label>
            <div class="aboutMe">{user.aboutMe}</div>
        {/if}
        <label for activity>Latest posts</label>
        {#if user.posts.length < 1}
            <Matrix />
            {:else}
        <div class="activity">
            {#each user.posts as post }
                <div>{post}</div>
            {/each}
        </div>
        {/if}
    </div>
</main>

<style>

main {
        display: flex;
        flex-direction: column;
        font-size: small;
    }

    img {
        max-width: 264px;
    }
    label {
        padding: 8px;
    }


    .aboutMe, .activity, .birthday{
        font-size: small;
        max-height: 200px;
        overflow: auto;
        border: solid 1px #333;
        border-radius: 6px;
        text-align: center;
        padding: 8px;
    }

    .activity {
        max-height: 500px;
    }
    

</style>