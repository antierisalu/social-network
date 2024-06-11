<script>
    import Matrix from '../../shared/matrix.svelte';
    import { slide } from 'svelte/transition';
    import { userInfo, userProfileData } from '../../stores'

    
    $userProfileData = $userInfo
    $: user = $userProfileData
    //  user.followers = ['DJ Worker Doctor', 'Doctor','DJ Worker Doctor', 'Producer DJ Worker','Producer DJ Worker', 'Doctor','DJ Worker Doctor', 'Doctor','DJ Worker Doctor', 'Producer DJ Worker','Producer DJ Worker', 'Doctor']
    // user.following = ['DJ Worker Doctor', 'Producer DJ Worker', 'Doctor']

</script>

<div class="PrivateData" in:slide out:slide>
    <label for="birthday">Birthday</label>
    <div class="birthday">{user.dateOfBirth.String}</div>
    {#if user.aboutMe.String}
        <label for="aboutMe">About me</label>
        <div class="aboutMe">{user.aboutMe.String}</div>
    {/if}
    <div class="follow">
        <div>
            {#if user.followers && user.followers.length > 0}

            <label for="followers">Followers</label>
            <div>
                {#each user.followers as follower}
                <div class="followers">{follower}</div>
                {/each}
            </div>
            {/if}

        </div>
        <div>
            {#if user.following && user.following.length > 0}
            <label for="followers">Following</label>
            <div >
                {#each user.following as following}
                <div class="following">{following}</div>
                {/each}
            </div>
            {/if}
        </div>
    </div>
    <div class="userPostLabels">
    <label for activity>Latest posts</label>
    <u>See all posts</u>
</div>
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

<style>

    label{
        padding: 8px;
        font-weight: bold;
    }
    .userPostLabels{
        display: flex;
        flex-direction: row;
        justify-content: space-evenly;
    }

    u{
        padding: 8px; 
        cursor: grabbing;
    }

    .follow {
        padding: 0;
        display: flex;
        justify-content: center;
        max-height: 200px;
        overflow-y: scroll;
        scrollbar-width: thin;
        scrollbar-color:  greenyellow #011;
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