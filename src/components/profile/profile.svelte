<script>
    import Button from "../../shared/button.svelte";
    import Matrix from '../../shared/matrix.svelte';
    import PrivateData from "./privateData.svelte";

    import { userInfo, userProfileData, isEditingProfile} from '../../stores'
    import { fade, slide } from 'svelte/transition';

    let followingUser
    let followRequested  
    
    $userProfileData = $userInfo
    $: user = $userProfileData

    function toggleProfile() {
    sendProfilePrivacyStatus()
    }

    let newNickname = '';
    let newAboutMe = '';

    function toggleEdit() {
        $isEditingProfile = !$isEditingProfile;
        if (!$isEditingProfile) {
            user.nickName.String = newNickname;
            user.aboutMe = newAboutMe;
        } else {
            newNickname = user.nickName.String;
            newAboutMe = user.aboutMe;
        }
    }


    function handleNicknameChange() {

    }


    async function sendFollow(action, target){
        console.log("sendfollow:",action, target)
        try {
        const response = await fetch("/followRequest", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ action: action, target: target })
            });
        }
        catch (error){
            console.error("Error sending follow request: ", error.message)
        }
}

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

        {#if user.nickName.String && !$isEditingProfile}
        <p in:fade>({user.nickName.String})</p>
        {:else if $isEditingProfile}
                <input in:fade class="editProfileText" type="text" bind:value={newNickname} on:input={handleNicknameChange} />
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
            {#if followingUser } 
                <Button id="unFollowBtn" on:click={()=> sendFollow(-1, user.id)} >unFollow</Button>
                {:else if !followingUser && followRequested}
                <Button id="unFollowBtn" on:click={()=> sendFollow(-1, user.id)} >Cancel request</Button>
                {:else }
                <Button type="secondary" w84={true} id="followBtn" on:click={()=> sendFollow(!user.privacy ? 1 : 0, user.id)}>Follow</Button>
            {/if}
            <Button type="secondary" inverse={true} w84={true} id="chatBtn">Chat</Button>
        </div>

        {:else}
            <div class="btnEditPrivate">
            {#if $userInfo.privacy}
                <div in:fade><br><Button type="secondary" inverse={true} on:click={toggleProfile}>Set Public</Button></div>
            {:else}
                <div in:fade><br><Button inverse={true} on:click={toggleProfile}>Set Private</Button></div>
            {/if}
            <Button type="secondary" inverse={!$isEditingProfile} id="editBtn" on:click={toggleEdit}>{$isEditingProfile ? 'Save Profile' : 'Edit Profile'}</Button>
            </div>
        {/if}
        
        {#if user.privacy == 0 || $userInfo.id == user.id}
        <PrivateData />
        {/if}
    </div>
</main>

<style>

    .btnEditPrivate {
        display: flex;
        justify-content: space-evenly;
        align-items: flex-end;
    }

main {
        display: flex;
        flex-direction: column;
        font-size: small;
    }

    img {
        max-width: 264px;
    }

    .name {
        padding: 8px;
    }

    .editProfileText {
        width: 100%;
        text-align: center;
        border-color: greenyellow;
    }

</style>