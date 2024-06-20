<script>
    import Button from "../../shared/button.svelte";
    import Matrix from '../../shared/matrix.svelte';
    import PrivateData from "./privateData.svelte";
    import ChangeImage from "../../shared/imagePreview.svelte"
    import { sendMessage } from "../../websocket.js"

    import { userInfo, userProfileData, isEditingProfile, newAboutMeStore,  uploadImageStore} from '../../stores'
    import { fade } from 'svelte/transition';

    let followRequested

    $userProfileData = $userInfo
    $: user = $userProfileData

    const toggleProfile = () => sendProfilePrivacyStatus()

    let newNickname = '';

    let uploadImage;
    uploadImageStore.subscribe(value => {
    uploadImage = value;
    });

    export function toggleEdit() {
        $isEditingProfile = !$isEditingProfile;
        if (!$isEditingProfile) {
            user.nickName.String = newNickname;
            user.aboutMe.String = $newAboutMeStore;
            saveProfileChanges();

        } else {
            newNickname = user.nickName.String;
            $newAboutMeStore = user.aboutMe.String;
        }
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
        let link = "follow_" + (user.id).toString()
        sendMessage(JSON.stringify({ type: "followRequestNotif", data: link }))

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

async function saveProfileChanges() {
        console.log("newNickname", newNickname, "newAboutMe:", $newAboutMeStore)
        uploadImage().catch(error => {console.error('Error uploading the image:', error); });
        const response = await fetch('/editProfile', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            nickName: newNickname,
            aboutMe: $newAboutMeStore,
            // Avatar as well
        })
    });


    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }
}

    </script>
<main>
    <div class="userContainer">
        <div class="name">{user.firstName} {user.lastName}</div>

        {#if user.nickName.String && !$isEditingProfile}
        <p in:fade>({user.nickName.String})</p>
        {:else if $isEditingProfile}
                <input in:fade class="editProfileText" type="text" bind:value={newNickname} />
        {/if}

        {#if user.avatar && !$isEditingProfile}
            <div class="avatar">
                <img src={user.avatar} border="0" alt="avatar" />
            </div>
        {:else if $isEditingProfile}
            <div><ChangeImage inputIDProp="changeAvatarImage" fakeInputText="Upload new Avatar" style="border-color: greenyellow; width:242px"/></div>
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
            <div class="btnEditPrivate">
            {#if !$isEditingProfile}
                {#if $userInfo.privacy}
                    <div in:fade><br><Button type="secondary" inverse={true} on:click={toggleProfile}>Set Public</Button></div>
                {:else}
                    <div in:fade><br><Button inverse={true} on:click={toggleProfile}>Set Private</Button></div>
                {/if}
            {:else}
                <div in:fade><Button type="primary" on:click={() => $isEditingProfile = false}>Cancel edit</Button></div>
            {/if}
            <Button type="secondary" inverse={!$isEditingProfile} id="editProfileBtn" on:click={toggleEdit}>{$isEditingProfile ? 'Save Profile' : 'Edit Profile'}</Button>
            </div>
        {/if}
        {#if user.privacy === 0 || $userInfo.id === user.id || user.isFollowing === true}
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
        padding:8px;
        /* margin: 0; */
    }

</style>