<script>

    import { onMount } from 'svelte';
    import Button from "../shared/button.svelte";
    import Matrix from '../shared/matrix.svelte';
    import { fade, slide } from 'svelte/transition';


    // For buttons to work

    let followingUser
    let followRequested
    let notPrivateProfile = true
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
        followers: '',
        following: '',
    };

    function toggleProfile() {
    notPrivateProfile = !notPrivateProfile;
    }

    onMount(async () => {
        try {
            const response = await fetch('http://localhost:8080/session'); // Replace with your actual endpoint
            const data = await response.json();
            user = data;
                    // user.posts = ['123', 'Hello'];
        // user.avatar = ''
        user.followers = ['DJ Worker Doctor', 'Doctor','DJ Worker Doctor', 'Producer DJ Worker','Producer DJ Worker', 'Doctor',]
        user.following = ['DJ Worker Doctor', 'Producer DJ Worker', 'Doctor',]
        } catch (error) {
            console.error('Error fetching user data:', error);
        }
    });



    
    </script>



<main>
    
    
    <div class="userContainer">
        <div class="name">{user.firstName} {user.lastName}</div>
        {#if user.nickName}
        <p>({user.nickName})</p>
        {/if}
        {#if user.avatar}
            <div class="avatar">
                <img src={user.avatar} border="0" alt="avatar" />
            </div>
        {:else}
            <Matrix /><br>
        {/if}
        {#if notOwnPage}
        <div class="buttons">
            {#if followingUser } 
                <Button id="unFollowBtn">unFollow</Button>
                {:else if !followingUser && followRequested}
                <Button id="unFollowBtn">Cancel request</Button>
                {:else }
                <Button type="secondary" w84={true} id="followBtn">Follow</Button>
            {/if}
            <Button type="secondary" inverse={true} w84={true} id="chatBtn">Chat</Button>
        </div>
        {/if}

        {#if notPrivateProfile}
        <Button type="secondary" inverse={true} on:click={toggleProfile}>Public</Button>
        <div class="PrivateData" transition:slide>
            <label for="birthday">Birthday</label>
            <div class="birthday">{user.dateOfBirth}</div>
            {#if user.aboutMe}
                <label for="aboutMe">About me</label>
                <div class="aboutMe">{user.aboutMe}</div>
            {/if}
            <div class="follow">
                <div>
                    <label for="followers">Followers</label>
                    <div>
                        {#each user.followers as follower}
                        <div class="followers">{follower}</div>
                        {/each}
                    </div>
                </div>
                <div>
                    <label for="followers">Following</label>
                    <div >
                        {#each user.following as following}
                        <div class="following">{following}</div>
                        {/each}
                    </div>
                </div>
            </div>
            <label for activity>Latest posts</label>
            {#if user.posts === undefined || user.posts.length < 1}
                <Matrix />
                {:else}
            <div class="activity">
                {#each user.posts as post }
                    <div>{post}</div>
                {/each}
            </div>
            {/if}
        </div>
        {:else}
        <div in:slide><Button inverse={true} on:click={toggleProfile}>Private</Button></div>
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
    label{
        padding: 8px;
        font-weight: bold;
    }

    .follow {
        padding: 0;
        display: flex;
        justify-content: center;
    }

    .aboutMe, .activity, .birthday, .following, .followers{
        font-size: small;
        max-height: 200px;
        overflow: auto;
        border: solid 1px #333;
        border-radius: 6px;
        text-align: center;
        padding: 8px;
    }

    .following, .followers {
        margin: 4px;
        padding: 4px auto;
    }

    .activity {
        max-height: 500px;
    }
    

</style>