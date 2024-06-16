<script>
    import Button from "../../shared/button.svelte";
    import Matrix from '../../shared/matrix.svelte';
    import PrivateData from "./privateData.svelte";

    import { userInfo, userProfileData } from '../../stores'
    import { fade } from 'svelte/transition';

    let followRequested

    $userProfileData = $userInfo
    $: user = $userProfileData

    function toggleProfile() {
    sendProfilePrivacyStatus()
    }

    async function sendFollow(action, target){
        console.log("sendfollow:",action, target)
        try {
        const response = await fetch("/api/followers", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ action: action, target: target })
            });

        let followStatus = await response.text()
        console.log(followStatus)
        if (action == 1) {
            user.isFollowing = true
            $userProfileData.followers.length++
        } else if (action == -1) {
            $userProfileData.followers.length--
            user.isFollowing = false

        }
        }
        catch (error){
            console.error("Error sending follow request: ", error.message)
        }
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
        {#if $userInfo.id != user.id}<!-- if the rendered user is not client -->
        <div class="buttons">
            {#if user.isFollowing }
                <Button id="unFollowBtn" on:click={()=> sendFollow(-1, user.id)} >unFollow</Button>
                {:else if !user.isFollowing && followRequested}
                <Button id="unFollowBtn" on:click={()=> sendFollow(-2, user.id)} >Cancel request</Button>
                {:else }
                <Button type="secondary" w84={true} id="followBtn" on:click={()=> sendFollow(!user.privacy ? 1 : 0, user.id)}>Follow</Button>
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
        {#if user.privacy === 0 || $userInfo.id === user.id || user.isFollowing === true}
        <PrivateData />
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

</style>