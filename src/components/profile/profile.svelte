<script>
    import Button from "../../shared/button.svelte";
    import Matrix from '../../shared/matrix.svelte';
    import { userInfo, userProfileData } from '../../stores'
    import { fade, slide } from 'svelte/transition';

    let followingUser 
    let followRequested  

    $userProfileData = $userInfo
    $: user = $userProfileData

    function toggleProfile() {
    sendProfilePrivacyStatus()
    }

//$userInfo.privacy == 0, siis saada true ehk tahan panna privateiks
//$userInfo.privacy == 1, siis saada false ehk tahan panna publicuks
    async function sendProfilePrivacyStatus() {
        try {
        const response = await fetch("/privacy", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ newPrivacy: !$userInfo.privacy })
        });

        if (!response.ok) {
            throw new Error("Network response was not ok: " + response.statusText);
        }

        const result = await response.json();//returns {newPrivacy: true/false}
        $userInfo.privacy = result.newPrivacy


    } catch (error) {
        console.error("Error sending profile privacy status:", error.message);
    }
}

    // user.followers = ['DJ Worker Doctor', 'Doctor','DJ Worker Doctor', 'Producer DJ Worker','Producer DJ Worker', 'Doctor','DJ Worker Doctor', 'Doctor','DJ Worker Doctor', 'Producer DJ Worker','Producer DJ Worker', 'Doctor',]
    // user.following = ['DJ Worker Doctor', 'Producer DJ Worker', 'Doctor',]

    </script>
<main>
    <div class="userContainer">
        <div class="name">{user.firstName} {user.lastName}</div>
        {#if user.nickName.String}
        <p>({user.nickName.String})</p>
        {/if}
        {#if user.avatar}
            <div class="avatar">
                <img src={user.avatar} border="0" alt="avatar" />
            </div>
        {:else}
            <Matrix /><br>
        {/if}
        {#if $userInfo.id != user.id}
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
        {:else}
            {#if $userInfo.privacy}
                <div in:fade><br><Button type="secondary" inverse={true} on:click={toggleProfile}>Set Public</Button></div>
            {:else}
                <div in:fade><br><Button inverse={true} on:click={toggleProfile}>Set Private</Button></div>
            {/if}
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