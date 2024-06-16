<script>
    import Matrix from '../../shared/matrix.svelte';
    import { slide } from 'svelte/transition';
    import { userInfo, userProfileData } from '../../stores'
    import Followers from './followers.svelte';

    let showOverlay = false;
    let overlayInfo = [];

    function followOverlay(n){
        showOverlay = true
        if (n === 1){
        overlayInfo = user.followers
        console.log(n, user.followers)
        } else {
        overlayInfo = user.following
            console.log(n, user.following)
        }
    }

    $userProfileData = $userInfo
    $: user = $userProfileData
    //  user.followers = ['DJ Worker Doctor', 'Doctor','DJ Worker Doctor', 'Producer DJ Worker','Producer DJ Worker', 'Doctor','DJ Worker Doctor', 'Doctor','DJ Worker Doctor', 'Producer DJ Worker','Producer DJ Worker', 'Doctor']
    // user.following = ['DJ Worker Doctor', 'Producer DJ Worker', 'Doctor']
    function toggleOverlay() {
    showOverlay = !showOverlay;
  }
</script>
{#if showOverlay && overlayInfo}
    <Followers on:close={toggleOverlay} followers={overlayInfo} />
{/if}
<div class="PrivateData" in:slide out:slide>
    <label for="birthday">Birthday</label>
    <div class="birthday">{user.dateOfBirth.String}</div>
    {#if user.aboutMe.String}
        <label for="aboutMe">About me</label>
        <div class="aboutMe">{user.aboutMe.String}</div>
    {/if}
    <div class="follow">
        <div>
            <label for="followers">Followers</label>
            <div>
                <div class="followers" on:click={()=>followOverlay(1)}>{user.followers ? user.followers.length : 0}</div>
            </div>
        </div>
        <div>
            <label for="followers">Following</label>
            <div>
                <div class="following" on:click={()=>followOverlay(0)}>{user.following ? user.following.length : 0}</div>
            </div>
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

    .following:hover, .followers:hover {
        color: white;
        border-color: white;
        cursor: pointer;
    }

    .activity {
        max-height: 500px;
    }

</style>