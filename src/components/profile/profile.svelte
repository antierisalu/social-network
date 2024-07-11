<script>
  import Button from "../../shared/button.svelte";
  import Matrix from "../../shared/matrix.svelte";
  import PrivateData from "./privateData.svelte";
  import ChangeImage from "../../shared/imagePreview.svelte";
  import { sendMessage } from "../../websocket.js";
  import { getPosts, selectUser } from "../../utils";
  import { onMount } from "svelte";

  import {
    userInfo,
    userProfileData,
    isEditingProfile,
    newAboutMeStore,
    uploadImageStore,
  } from "../../stores";
  import { fade } from "svelte/transition";

  $: followerCount = $userProfileData.followers
    ? $userProfileData.followers.length
    : 0;
  const toggleProfile = () => sendProfilePrivacyStatus();

  let newNickname = "";

  let uploadImage;
  uploadImageStore.subscribe((value) => {
    uploadImage = value;
  });

  export function toggleEdit() {
    $isEditingProfile = !$isEditingProfile;
    if (!$isEditingProfile) {
      $userProfileData.nickName.String = newNickname;
      $userProfileData.aboutMe.String = $newAboutMeStore;
      saveProfileChanges();
    } else {
      newNickname = $userProfileData.nickName.String;
      $newAboutMeStore = $userProfileData.aboutMe.String;
    }
  }

  async function sendFollow(action, target) {
    console.log("sendfollow:", action, target);
    try {
      const response = await fetch("/api/followers", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ action: action, target: target }),
      });

      let userData = await response.json(); //returns who initiated follow change
      var messageData = {
        type: "follow",
        targetid: $userProfileData.id,
        fromid: $userInfo.id,
        data: String,
      };
      if (action == 0) {
        messageData.type = "followRequest";
      } else if (action === -2) {
        messageData.type = "cancelRequest";
      }

      //update frontend follower list
      if (userData.followStatus == 1) {
        $userProfileData.isFollowing = 1;
        $userProfileData.followers = $userProfileData.followers //add user to followers list, if followerslist is null make a new array
          ? [...$userProfileData.followers, userData.user]
          : [userData.user];

        //send notification
        messageData.data = "follow_" + messageData.fromid.toString();
      } else if (userData.followStatus == -1) {
        $userProfileData.isFollowing = -1;
        const objString = JSON.stringify(userData.user); //remove user from followers list
        $userProfileData.followers = $userProfileData.followers.filter(
          (item) => JSON.stringify(item) !== objString,
        );
      } else if (userData.followStatus == 0) {
        messageData.data = "followRequest_" + messageData.fromid.toString();
      } else if (userData.followStatus == -2) {
        messageData.data = "cancelRequest_" + messageData.fromid.toString();
      }
      sendMessage(JSON.stringify(messageData));
    } catch (error) {
      console.error("Error sending follow request: ", error.message);
    }
    getPosts();
    console.log($userProfileData);
    selectUser($userProfileData.id); //Reload profile to reset allposts, followers, etc.
  }

  async function sendProfilePrivacyStatus() {
    try {
      const response = await fetch("/privacy", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ newPrivacy: !$userInfo.privacy }),
      });

      if (!response.ok) {
        throw new Error("Network response was not ok: " + response.statusText);
      }

      const result = await response.json(); //returns {newPrivacy: true/false}
      $userInfo.privacy = result.newPrivacy;
    } catch (error) {
      console.error("Error sending profile privacy status:", error.message);
    }
  }

  async function saveProfileChanges() {
    let path = await uploadImage().catch((error) => {
      console.error("Error uploading the image:", error);
    });
    if (path === undefined) {
      path = $userInfo.avatar;
    }
    const response = await fetch("/editProfile", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        nickName: newNickname,
        aboutMe: $newAboutMeStore,
        avatar: path,
      }),
    });
    console.log(path);
    if (path !== undefined) {
      $userInfo.avatar = path;
      $userProfileData = $userInfo;
    }
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
  }
</script>

<main>
  <div class="userContainer">
    <div class="name">
      {$userProfileData.firstName}
      {$userProfileData.lastName}
    </div>

    {#if $userProfileData.nickName.String && !$isEditingProfile}
      <p in:fade>({$userProfileData.nickName.String})</p>
    {:else if $isEditingProfile}
      <input
        in:fade
        class="editProfileText"
        type="text"
        placeholder="Add Nickname"
        bind:value={newNickname}
      />
    {/if}

    {#if $userProfileData.avatar && !$isEditingProfile}
      <div class="avatar">
        <img src={$userProfileData.avatar} alt="avatar" />
      </div>
    {:else if $isEditingProfile}
      <div>
        <ChangeImage
          src={$userInfo.avatar}
          inputIDProp="changeAvatarImage"
          fakeInputText="Upload new Avatar"
          style="border-color: greenyellow; width:242px"
        />
      </div>
    {:else}
      <Matrix /><br />
    {/if}

    {#if $userInfo.id != $userProfileData.id}<!-- if the rendered user is not client -->
      <div class="buttons">
        {#if $userProfileData.areFollowing == 1}<!-- 1 == am following -->
          <Button
            id="unFollowBtn"
            on:click={() => sendFollow(-1, $userProfileData.id)}
            >unFollow</Button
          >
        {:else if $userProfileData.areFollowing == 0}
          <!-- 0 == i have requested -->
          <Button id="unFollowBtn" on:click={() => sendFollow(-2, $userProfileData.id)}
            >Cancel request</Button
          >
        {:else}
          <Button
            type="secondary"
            w84={true}
            id="followBtn"
            on:click={() =>
              sendFollow(
                !$userProfileData.privacy ? 1 : 0,
                $userProfileData.id,
              )}>Follow</Button
          >
        {/if}
        <Button type="secondary" inverse={true} w84={true} id="chatBtn"
          >Chat</Button
        >
      </div>
    {:else}
      <div class="btnEditPrivate">
        {#if !$isEditingProfile}
          {#if $userInfo.privacy}
            <div in:fade>
              <br /><Button
                type="secondary"
                inverse={true}
                on:click={toggleProfile}>Set Public</Button
              >
            </div>
          {:else}
            <div in:fade>
              <br /><Button inverse={true} on:click={toggleProfile}
                >Set Private</Button
              >
            </div>
          {/if}
        {:else}
          <div in:fade>
            <Button type="primary" on:click={() => ($isEditingProfile = false)}
              >Cancel edit</Button
            >
          </div>
        {/if}
        <Button
          type="secondary"
          inverse={!$isEditingProfile}
          id="editProfileBtn"
          on:click={toggleEdit}
          >{$isEditingProfile ? "Save Profile" : "Edit Profile"}</Button
        >
      </div>
    {/if}
    {#if $userProfileData.privacy === 0 || $userInfo.id === $userProfileData.id || $userProfileData.areFollowing === 1}
      <PrivateData {followerCount} />
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
    width: 200px;
    border-radius: 20px;
  }

  .name {
    padding: 8px;
  }

  .editProfileText {
    width: 100%;
    text-align: center;
    border-color: greenyellow;
    padding: 8px;
    /* margin: 0; */
  }
</style>
